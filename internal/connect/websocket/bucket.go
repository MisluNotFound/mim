package websocket

import "sync"

type Bucket struct {
	clients map[int64]*Client
	lock    sync.RWMutex
}

func (b *Bucket) RemoveClient(id int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.clients, id)
}

func NewBucket() (b *Bucket) {
	b = new(Bucket)
	b.clients = make(map[int64]*Client, 50)
	return
}

func (b *Bucket) LRU() {

}