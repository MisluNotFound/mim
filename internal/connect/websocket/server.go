/*
server存储在线用户列表，活跃的聊天室
*/
package websocket

import (
	"crypto/md5"
	logicrpc "mim/internal/connect/rpc/logic_rpc"
	"mim/pkg/proto"
	"strconv"
)

var Default *Server

type Server struct {
	Bucket    []*Bucket
	ServerID  int
	Count     int
}

func NewServer(buckets []*Bucket, id int) *Server {
	return &Server{
		ServerID: id,
		Bucket:   buckets,
		Count:    len(buckets),
	}
}

func (s *Server) getHashCode(id int64) int {
	idStr := strconv.FormatInt(id, 10)
	h := md5.New()
	h.Write([]byte(idStr))
	hashBytes := h.Sum(nil)
	hashValue := int(hashBytes[0] | hashBytes[1]<<8 | hashBytes[2]<<16 | hashBytes[3]<<24)
	bucketIdx := hashValue % s.Count
	return bucketIdx
}

func (s *Server) assignUser(c *Client) {
	bucketIdx := s.getHashCode(c.ID)

	_, ok := s.GetUser(c.ID)
	if ok {
		return
	}

	b := s.Bucket[bucketIdx]
	b.lock.Lock()
	defer b.lock.Unlock()

	b.clients[c.ID] = c
}

func (s *Server) getBucket(id int64) *Bucket {
	bucketIdx := s.getHashCode(id)

	return s.Bucket[bucketIdx]
}

func (s *Server) GetUser(id int64) (*Client, bool) {
	b := s.getBucket(id)
	b.lock.RLock()
	defer b.lock.RUnlock()

	c, ok := b.clients[id]
	return c, ok
}

func (s *Server) AssignInBucket(c *Client) {

	req := &proto.OnlineReq{
		UserID:   c.ID,
		ServerID: s.ServerID,
	}
	if err := logicrpc.Online(req); err != nil {
		c.lock.Lock()
		c.Conn.WriteJSON("登录失败，请重新登录")
		c.lock.Unlock()
		return
	}
	s.assignUser(c)
}
