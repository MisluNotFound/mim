package proto

import (
	"mim/internal/logic/dao"
	"mim/pkg/code"
)

type MessageReq struct {
	SenderID int64
	TargetID int64
	Ack      int64
	Type     int
	Body     []byte
}

type MessageResp struct {
}

type PushMessageReq struct {
	SenderID int64
	TargetID int64
	Seq      int64
	Body     []byte
}

type PushMessageResp struct {
	IsOffline bool
}

type PullMessageReq struct {
	UserID   int64
	TargetID int64
	LastSeq  int
	Size     int
}

type PullMessageResp struct {
	Code     code.ResCode
	Messages map[int64][]dao.Message
}
