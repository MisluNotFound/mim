// 根据消息类型选择不同的函数
package messaging

import (
	"fmt"
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

func singleHandler(msg *redis.Message) {
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
		Body:     msg.Body,
	}
	fmt.Println(info)
	if info.UserID == 0 {
		msg.Status = StatusOffline
	} else {
		msg.Status = StatusOnline
		wsrpc.PushMessage(req)
	}

	zap.L().Info("singleHandler receive message", zap.Any("msg", req))
	go asyncSaveMessage(msg)
}

func groupHandler(msg *redis.Message) {

}

func pongHandler(msg *redis.Message) {

}

func ackHandler(msg *redis.Message) {

}

func asyncSaveMessage(msg *redis.Message) {
	redis.StoreRedisMessage(*msg)
	dao.StoreMysqlMessage(&dao.Message{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Content:  msg.Body,
		Seq:      msg.Seq,
	})
}
