package snowflake

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

var mu sync.Mutex

func GenerateUniqueID() int64 {
	mu.Lock()
	defer mu.Unlock()

	// 获取当前时间戳（精度到毫秒）
	timestamp := time.Now().UnixNano() / 1e6

	randNum, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}

	uniqueID := timestamp*1000000 + randNum.Int64()

	return uniqueID
}
