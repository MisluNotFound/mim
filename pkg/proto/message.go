package proto

import (
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"time"
)

type MessageReq struct {
	SenderID int64
	TargetID int64
	Ack      int64
	Type     int
	Body     []byte
	Media    string
	URL      string
}

type MessageResp struct {
}

type PushMessageReq struct {
	SenderID int64
	TargetID int64
	Seq      int64
	Body     []byte
	Timer    time.Time
	Extra    interface{}
}

type PushMessageResp struct {
	IsOffline bool
}

type PullMessageReq struct {
	UserID    int64
	SessionID int64
	LastSeq   int64
	Size      int
	IsGroup   bool
}

type PullMessageResp struct {
	Code     code.ResCode
	Messages []dao.Message
}

type OfflineMessageReq struct {
	SenderID int64
	TargetID int64
	Seq      int64
	Body     []byte
}

type PullOfflineMessageReq struct {
	UserID    int64
	SessionID int64
	Count     int
	IsGroup   bool
}

type GetUnReadCountReq struct {
	UserID int64
}

type UnReadInfo struct {
	SessionID   int64
	Count       int
	LastMessage dao.Message
}

type GetUnReadResp struct {
	Code        code.ResCode
	SessionInfo []UnReadInfo
}

type PullErrMessageReq struct {
	UserID int64
}

type PullErrMessageResp struct {
	Code     code.ResCode
	Messages []dao.Message
}
