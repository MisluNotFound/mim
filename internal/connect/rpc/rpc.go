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
	zap.L().Info("init connect rpc server success")
	s.Serve("tcp", "8084")
}

