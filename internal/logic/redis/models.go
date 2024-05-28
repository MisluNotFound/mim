package redis

type UserInfo struct {
	UserID   int64
	ServerID int
	BucketID int
	IsNotice int // 用于查看是否通知用户有离线消息
}

type Message struct {
	Seq      int64 //id
	SenderID int64
	TargetID int64
	Body     []byte
	Status   string
}
