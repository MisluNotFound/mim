package websocket

import (
	"container/heap"
	"crypto/md5"
	"encoding/binary"
	logicrpc "mim/internal/connect/logic_rpc"
	"mim/pkg/mq"
	"mim/pkg/proto"
	"strconv"
	"time"
)

var Default *Server

type Server struct {
	Bucket    []*Bucket
	ServerID  int
	Count     int
	RabbitMQ  *mq.RabbitMQ
	publisher *mq.Publisher
}

func NewServer(bucketSize, id int, rabbitMQURL string) *Server {
	rabbitMQ, _ := mq.NewRabbitMQ(rabbitMQURL)
	server := &Server{
		ServerID:  id,
		RabbitMQ:  rabbitMQ,
		publisher: mq.NewPublisher(rabbitMQ),
	}

	buckets := make([]*Bucket, bucketSize)
	for i := range buckets {
		buckets[i] = NewBucket(rabbitMQ, server.ServerID, i)
	}

	server.Bucket = buckets
	server.Count = len(buckets)

	for i := range server.Bucket {
		go server.Bucket[i].Consumer.Work(false, consumeMessage)
	}

	return server
}

func (s *Server) getHashCode(id int64) int {
	idStr := strconv.FormatInt(id, 10)
	h := md5.New()
	h.Write([]byte(idStr))
	hashBytes := h.Sum(nil)

	hashValue := int(binary.BigEndian.Uint32(hashBytes[:4]))

	bucketIdx := hashValue % s.Count
	return bucketIdx
}

func (s *Server) assignUser(c *Client) {
	bucketIdx := s.getHashCode(c.ID)
	c.BucketID = bucketIdx

	_, ok := s.GetUser(c.ID)
	if ok {
		return
	}

	b := s.Bucket[bucketIdx]
	if b.capacity == b.cHeap.Len() {
		b.LRU()
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	c.HeartBeat = time.Now()

	heap.Push(&b.cHeap, c)

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
	if !ok {
		return nil, false
	}

	return c, true
}

func (s *Server) AssignInBucket(c *Client) {

	req := &proto.OnlineReq{
		UserID:   c.ID,
		ServerID: s.ServerID,
		BucketID: c.BucketID,
	}

	if err := logicrpc.Online(req); err != nil {
		c.lock.Lock()
		c.Conn.WriteJSON("登录失败，请重新登录")
		c.lock.Unlock()
		return
	}

	s.assignUser(c)
}

func (s *Server) Close() {

}
