package websocket

import "sync"

type Room struct {
	ID    int64
	Count int
	lock  sync.RWMutex
	head  *Client
}

func (r *Room) AddClient(c *Client) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.head == nil {
		r.head = c
		return
	}

	r.head.pre = c
	c.next = r.head
	r.Count++
}
