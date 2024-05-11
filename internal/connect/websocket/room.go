package websocket

import "sync"

type Room struct {
	ID    int64
	Count int
	lock  sync.RWMutex
	next  *Client
}
