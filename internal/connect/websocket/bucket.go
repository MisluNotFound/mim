package websocket

import (
	"mim/setting"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Bucket struct {
	pool    []*Client // 连接池
	clients map[int64]*Client
	rooms   map[int64]*Room
	lock    sync.RWMutex
	cond    *sync.Cond
	maxSize int
}

func NewBucket(maxSize int) *Bucket {
	b := &Bucket{
		clients: make(map[int64]*Client),
		maxSize: maxSize,
		pool:    make([]*Client, maxSize),
	}
	b.cond = sync.NewCond(&b.lock)
	return b
}

// 向连接池里添加连接
func (b *Bucket) addConnect(conn *Client) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.pool) < b.maxSize {
		b.pool = append(b.pool, conn)
		b.cond.Signal()
	} else {
		conn.Conn.Close()
	}
}

func (b *Bucket) getConnect() *Client {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.pool) == 0 {
		b.cond.Wait()
		// 应当启用用户检测的算法释放连接
	}

	conn := b.pool[0]
	b.pool = b.pool[1:]
	return conn
}

func (b *Bucket) releaseConnection(client *Client) {
	client.lock.Lock()
	client.IsUse = false
	defer client.lock.Unlock()

	b.addConnect(client)
}

func newClientPool(b *Bucket, url string) {
	for i := 0; i < b.maxSize; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			zap.L().Error("init client pool failed: ", zap.Error(err))
			continue
		}
		b.addConnect(conn)
	}
}

// 添加到在线表中
func (b *Bucket) AddClient(id int64, username string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	c := NewClient(id, username, setting.Conf.WsConfig.ChannelSize)
	conn := 1
}

func (b *Bucket) GetClient() {

}

func (b *Bucket) RemoveClient() {

}
