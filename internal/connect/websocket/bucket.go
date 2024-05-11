package websocket

import "sync"

type Bucket struct {
	clients map[int64]*Client
	lock    sync.RWMutex
	rooms   map[int64]*Room
}

func (b *Bucket) RemoveClient(id int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.clients, id)
}

func NewBucket() (b *Bucket) {
	b = new(Bucket)
	b.clients = make(map[int64]*Client, 50)
	b.rooms = make(map[int64]*Room, 20)
	return
}

func (b *Bucket) LRU() {

}
