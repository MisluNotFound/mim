package proto

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
	Req      int64
	Body     []byte
}

type PushMessageResp struct {
	IsOffline bool
}
