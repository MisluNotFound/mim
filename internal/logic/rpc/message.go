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

// 获取历史记录
func (r *LogicRpc) PullMessage(ctx context.Context, req *proto.PullMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess
	fmt.Println("userID", req.UserID, "seq", req.LastSeq, "sessionID", req.SessionID)
	// 判断是否为群
	if req.IsGroup {
		// 获取入群时间
		ug, ok, err := dao.IsJoined(req.UserID, req.SessionID)
		if err != nil {
			resp.Code = code.CodeServerBusy
			zap.L().Error("pull message failed: ", zap.Error(err))
			return err
		}

		if !ok {
			resp.Code = code.CodeNotJoinGroup
			return dao.ErrorNotJoinGroup
		}

		resp.Messages, err = dao.PullGroupMessage(req.SessionID, req.LastSeq, ug.JoinTime, req.Size)
		if err != nil {
			resp.Code = code.CodeServerBusy
			zap.L().Error("pull message failed: ", zap.Error(err))
			return err
		}
	} else {
		var err error
		resp.Messages, err = dao.PullSingleMessage(req.UserID, req.SessionID, req.LastSeq, req.Size)
		if err != nil {
			resp.Code = code.CodeServerBusy
			zap.L().Error("pull message failed: ", zap.Error(err))
			return err
		}
	}
	return nil
}

// 用户点开会话获取未读消息
func (r *LogicRpc) PullOfflineMessage(ctx context.Context, req *proto.PullOfflineMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess
	var messages []dao.Message
	var err error
	if req.IsGroup {
		var ok bool
		var ug *dao.UserGroup
		ug, ok, err = dao.IsJoined(req.UserID, req.SessionID)
		if err != nil {
			zap.L().Error("pull offline message failed: ", zap.Error(err))
			resp.Code = code.CodeServerBusy
			return err
		}

		if !ok {
			resp.Code = code.CodeNotJoinGroup
			return nil
		}

		fmt.Println("user group: ", ug)
		messages, err = dao.PullOfflineMessage(req.UserID, req.SessionID, true, req.Count, ug.JoinTime)
	} else {
		messages, err = dao.PullOfflineMessage(req.UserID, req.SessionID, true, req.Count, 0)
	}
	fmt.Println(req)

	fmt.Println(messages)
	if err != nil {
		zap.L().Error("pull offline message failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	err = redis.MarkAsRead(req.UserID, req.SessionID)
	if err != nil {
		zap.L().Error("pull offline message failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	resp.Messages = messages
	return nil
}

// 获取会未读消息数
func (r *LogicRpc) GetUnReadCount(ctx context.Context, req *proto.GetUnReadCountReq, resp *proto.GetUnReadResp) error {
	resp.Code = code.CodeSuccess

	// redis获取count
	counts, err := redis.GetUnReadCount(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	// db获取最后一条消息
	var senders []int64
	for s := range counts {
		senders = append(senders, s)
	}
	lastMessages, err := dao.GetLastMessage(senders, req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}
	// 构造返回值
	var infos []proto.UnReadInfo
	for s, c := range counts {
		info := proto.UnReadInfo{
			SessionID: s,
			Count:     c,
		}
		for _, m := range lastMessages {
			if m.SenderID == s || m.TargetID == s {
				info.LastMessage = m
			}
		}
		infos = append(infos, info)
	}
	resp.SessionInfo = infos
	return nil
}

func (r *LogicRpc) PullErrMessage(ctx context.Context, req *proto.PullErrMessageReq, resp *proto.PullErrMessageResp) error {
	resp.Code = code.CodeSuccess

	// 获取lastAck
	lastAck, err := redis.GetLastAck(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if lastAck == -1 {
		return nil
	}

	// 查出lastAck之后的消息
	messages, err := dao.PullErrMessage(req.UserID, lastAck)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	resp.Messages = messages
	return nil
}
