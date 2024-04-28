package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineId int64) error {
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}

	snowflake.Epoch = st.UnixNano()
	node, err = snowflake.NewNode(machineId)
	return err
}

func GenID() int64 {
	return node.Generate().Int64()
}