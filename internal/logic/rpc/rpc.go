package rpc

import (
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

type LogicRpc struct {
}

func InitLogicRpc() {
	s := server.NewServer()
	if err := s.RegisterName("LogicRpc", new(LogicRpc), ""); err != nil {
		zap.L().Error("init logicRpc failed: ", zap.Error(err))
	}
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})

	zap.L().Info("init logicRpc success")
	s.Serve("tcp", "localhost:8081")
}
