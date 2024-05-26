package websocket

import (
	logicrpc "mim/internal/connect/logic_rpc"
	"mim/pkg/proto"
	"mim/setting"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader websocket.Upgrader

func InitWebsocket() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  setting.Conf.WsConfig.ReadBufferSize,
		WriteBufferSize: setting.Conf.WsConfig.WriteBufferSize,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serve(Default, w, r)
	})
	zap.L().Info("init ws server success")
	if err := http.ListenAndServe(setting.Conf.WsConfig.WSServers[0].Addr, nil); err != nil {
		zap.L().Error("init ws server failed: ", zap.Error(err))
	}
}

func serve(s *Server, w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.L().Error("server() failed: ", zap.Error(err))
		return
	}

	req := &proto.AuthReq{
		Token: token,
	}

	id, username, err := logicrpc.Auth(req)
	if err != nil {
		zap.L().Error("unauthorized")
		conn.Close()
		return
	}

	c := NewClient(id, username, setting.Conf.WsConfig.ChannelSize)

	c.Conn = conn
	c.server = s
	c.ID = id

	// client放入在线表
	handleConnection(s, c)
}

func handleConnection(s *Server, c *Client) {
	s.AssignInBucket(c)
	go c.readProc()
	go c.writeProc()
	zap.L().Info("user connected: ", zap.Any("id", c.ID))
}
