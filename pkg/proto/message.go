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
	LastSeq  int64
	Size     int
}

type Data struct {
	Sessions []string
	Messages map[string][]dao.Message
}

type PullMessageResp struct {
	Code code.ResCode
	Data Data
}

type OfflineMessageReq struct {
	SenderID int64
	TargetID int64
	Seq      int64
	Body     []byte
}
