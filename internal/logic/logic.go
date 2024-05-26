package logic

import (
	wsrpc "mim/internal/logic/pconn_rpc/conn_rpc"
	"mim/internal/logic/pconn_rpc/messaging/mq"
	prpc "mim/internal/logic/pconn_rpc/rpc"
	"mim/internal/logic/redis"
	"mim/internal/logic/rpc"
	"mim/setting"
)

func InitLogic() {
	go rpc.InitLogicRpc()
	go prpc.InitLogicRpc()
	go wsrpc.InitWsRpc()
	redis.Close()
	go mq.InitMQ(setting.Conf.MQConfig.URL, setting.Conf.Exchange,
		setting.Conf.MQConfig.Queue, setting.Conf.MQConfig.RoutingKey,
		setting.Conf.MQConfig.LogicConsumersNum, setting.Conf.MQConfig.LogicPublishersNum)
}
