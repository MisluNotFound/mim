/*
server存储在线用户列表，活跃的聊天室
提供读写协程
*/
package websocket

import (
	"crypto/md5"
	"strconv"
)

var Default *Server

type Server struct {
	Bucket []*Bucket 
	Count int
}

func NewServer(buckets []*Bucket) *Server {
	return &Server{
		Bucket: buckets,
		Count: len(buckets),
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

func (s *Server) AssignToBucket(c *Client) {
	bucketIdx := s.getHashCode(c.ID)

	b := s.Bucket[bucketIdx]
	b.lock.Lock()
	defer b.lock.Unlock()

	b.clients[c.ID] = c
}

func (s *Server) getUserBucket(id int64) *Bucket {
	bucketIdx := s.getHashCode(id)

	return s.Bucket[bucketIdx]
} 

func (s *Server) GetUser(id int64) *Client {
	b := s.getUserBucket(id)
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.clients[id]
}