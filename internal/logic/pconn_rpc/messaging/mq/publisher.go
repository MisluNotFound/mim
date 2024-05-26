package mq

import (
	"mim/pkg/mq"
	"sync/atomic"
)

var (
	publishers     []*mq.Publisher
	publisherIndex int32
)

func GetPublisher() *mq.Publisher {
	idx := atomic.AddInt32(&publisherIndex, 1)
	return publishers[idx%int32(len(publishers))]
}
