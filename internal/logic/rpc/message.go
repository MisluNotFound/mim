package rpc

import (
	"context"
	"fmt"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/code"
	"mim/pkg/proto"
)

func (r *LogicRpc) PullMessage(ctx context.Context, req *proto.PullMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess

	// redis获取消息seq数组
	seqs, err := redis.GetMessages(req.UserID, req.TargetID, req.LastSeq, req.Size)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}
	fmt.Println("online messages", seqs)

	offlineSeqs, _ := redis.GetOfflineMessages(req.UserID)
	fmt.Println("offline messages", seqs)
	seqs = append(seqs, offlineSeqs...)

	// mysql获取消息内容
	messages, err := dao.GetMessages(seqs)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	// 包装 map[session][]messages
	data := map[int64][]dao.Message{}
	for _, msg := range messages {
		s, t := msg.SenderID, msg.TargetID

		if s != req.UserID {
			data[s] = append(data[s], msg)
		} else if t != req.UserID {
			data[t] = append(data[t], msg)
		}
	}

	resp.Messages = data
	return nil
}
