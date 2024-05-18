// prpc为长连接提供服务
package prpc

import (
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

type PRpc struct {
}

func InitLogicRpc() {
	s := server.NewServer()
	if err := s.RegisterName("PRpc", new(PRpc), ""); err != nil {
		zap.L().Error("init PRpc failed: ", zap.Error(err))
	}
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})

	zap.L().Info("init PRpc success")
	s.Serve("tcp", "localhost:8083")
}

