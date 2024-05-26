// 根据消息类型选择不同的函数
package messaging

import (
	"mim/internal/logic/dao"

	wsrpc "mim/internal/logic/pconn_rpc/conn_rpc"
	"mim/internal/logic/redis"
	"mim/pkg/proto"
	"mim/pkg/snowflake"

	"go.uber.org/zap"
)

var (
	StatusOnline  = "online"
	StatusOffline = "offline"
)

func singleHandler(msg *dao.Message) {
	// 查询用户状态
	info, err := redis.GetUserInfo(msg.TargetID)
	if err != nil {
		zap.L().Error("singleHandler() failed: ", zap.Error(err))
	}
	msg.Seq = snowflake.GenID()

	req := &proto.PushMessageReq{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Seq:      msg.Seq,
		Body:     msg.Content,
	}
	var status string
	if info.UserID == 0 {
		status = StatusOffline
	} else {
		status = StatusOnline
		go wsrpc.PushMessage(req)
	}

	// zap.L().Info("singleHandler receive message", zap.Any("msg", req))
	go asyncSaveMessage(msg, status)
}

func groupHandler(msg *redis.Message) {

}

func pongHandler(msg *redis.Message) {

}

func ackHandler(msg *redis.Message) {

}

func asyncSaveMessage(msg *dao.Message, status string) {
	redis.StoreRedisMessage(*msg, status)
	dao.StoreMysqlMessage(&dao.Message{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Content:  msg.Content,
		Seq:      msg.Seq,
	})
}
