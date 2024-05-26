package proto

import "mim/pkg/code"

type OnlineReq struct {
	UserID   int64
	ServerID int
	BucketID int
}

type OnlineResp struct {
	Code code.ResCode
}

type OfflineReq struct {
	UserID int64
}

type OfflineResp struct {
	Code code.ResCode
}
