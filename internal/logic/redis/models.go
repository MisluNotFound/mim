package redis

type UserInfo struct {
	UserID   int64 
	ServerID int  
}

type Message struct {
	Seq      int64 //id
	SenderID int64
	TargetID int64
	Body     []byte
	Status   string
}


