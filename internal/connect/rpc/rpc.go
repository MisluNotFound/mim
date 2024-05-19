package rpc

import (
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

type ConnectRpc struct {
}

func InitConnectRpc() {
	s := server.NewServer()
	if err := s.RegisterName("ConnectRpc", new(ConnectRpc), ""); err != nil {
		zap.L().Error("init connect rpc server failed: ", zap.Error(err))
		return 
	}
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})
	if err := s.Serve("tcp", "localhost:8084"); err != nil {
		zap.L().Error("init connect rpc server failed: ", zap.Error(err))
		return
	}

	zap.L().Info("init connect rpc server success")
}

