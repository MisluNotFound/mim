package rpc

import (
	"context"
	"fmt"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func (r *LogicRpc) PullMessage(ctx context.Context, req *proto.PullMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess
	data := map[string][]dao.Message{}

	// 获取用户会话
	userSessions, err := redis.GetUserSessions(req.UserID)
	if err != nil {
		zap.L().Error("get userSession failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	// 尝试从缓存中获取
	messages, err := redis.GetMessages(req.UserID, req.TargetID, req.LastSeq, req.Size)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	fmt.Println("online messages", messages)
	// 未命中
	if len(messages) < req.Size {
		zap.L().Info("cache miss, get message from mysql")
		messages, err = dao.GetMessages(req.UserID, req.TargetID, req.LastSeq, req.Size)
		if err != nil {
			resp.Code = code.CodeServerBusy
			return err
		}
		fmt.Println("read message from mysql", messages)
		go redis.WriteBack(req.UserID, messages)
	}
	
	// 包装 map[session][]messages
	for _, msg := range messages {
		session := redis.GetSessionID(msg.SenderID, msg.TargetID)
		data[session] = append(data[session], msg)
	}
	
	offlineMessages, _ := redis.GetOfflineMessages(req.UserID)
	fmt.Println("offline messages", offlineMessages)

	// 可以优化 慢
	for _, msg := range offlineMessages {
		session := redis.GetSessionID(msg.SenderID, msg.TargetID)
		for i := range data[session] {

			if data[session][i].Seq == msg.Seq {
				data[session][i].IsRead = true
			}
		}
	}

	resp.Data.Sessions = userSessions
	resp.Data.Messages = data
	return nil
}