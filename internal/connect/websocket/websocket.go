package websocket

import (
	"mim/internal/connect/rpc"
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
	if err := http.ListenAndServe(setting.Conf.WsConfig.Addr, nil); err != nil {
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

	id, username, err := rpc.Auth(req)
	if err != nil {
		zap.L().Error("unauthorized")
		conn.Close()
		return
	}

	c := NewClient(setting.Conf.WsConfig.ChannelSize)
	c.Conn = conn
	c.server = s
	c.ID = id
	c.Username = username
	c.done = make(chan struct{})
	c.channel = make(chan []byte, 20)
	s.AssignToBucket(c)
	go c.readProc()
	go c.writeProc()
	zap.L().Info("user connected: ", zap.Any("id", id))
}
