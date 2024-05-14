/*
server存储在线用户列表，活跃的聊天室
*/
package websocket

import (
	"crypto/md5"
	"mim/internal/connect/rpc"
	"mim/pkg/proto"
	"strconv"

	"github.com/gorilla/websocket"
)

var Default *Server

type Server struct {
	Bucket []*Bucket
	Count  int
}

func NewServer(buckets []*Bucket) *Server {
	return &Server{
		Bucket: buckets,
		Count:  len(buckets),
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
	defer b.lock.Unlock()

	c, ok := b.clients[id]
	return c, ok
}

func (s *Server) GetRoom(id int64) (*Room, bool) {
	b := s.getBucket(id)
	b.lock.RLock()
	defer b.lock.Unlock()

	r, ok := b.rooms[id]
	return r, ok
}

func (s *Server) assignRoom(g *Room) {
	bucketIdx := s.getHashCode(g.ID)

	_, ok := s.GetRoom(g.ID)
	if ok {
		return
	}

	b := s.Bucket[bucketIdx]
	b.lock.Lock()
	defer b.lock.Unlock()

	b.rooms[g.ID] = g
}

func (s *Server) AssignInBucket(c *Client) {
	groups, err := rpc.GetGroup(&proto.FindGroupsReq{
		UserID: c.ID,
	})

	s.assignUser(c)
	if err != nil {
		c.Conn.WriteMessage(websocket.BinaryMessage, []byte("server busy"))
		return
	}

	for _, g := range *groups {
		r := &Room{
			ID:    g,
			Count: 0,
		}
		s.assignRoom(r)
		r.AddClient(c)
	}
}
