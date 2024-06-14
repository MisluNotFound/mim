package logic

import (
	"mim/internal/logic/mq"
	"mim/internal/logic/redis"
	"mim/internal/logic/rpc"
	"mim/setting"
)

func InitLogic() {
	go rpc.InitLogicRpc()
	redis.Close()
	go mq.InitMQ(setting.Conf.MQConfig.URL, setting.Conf.Exchange,
		setting.Conf.MQConfig.Queue, setting.Conf.MQConfig.RoutingKey,
		setting.Conf.MQConfig.LogicConsumersNum, setting.Conf.MQConfig.LogicPublishersNum)
}
