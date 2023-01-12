package id

import (
	"math/rand"
	"sync/atomic"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	sn, err := snowflake.NewNode(rand.Int63n(1024))
	if err != nil {
		panic(err)
	}
	node = sn
}
func GenID() int64 {
	return node.Generate().Int64()
}

var inc int32

func GenTempID() int32 {
	return atomic.AddInt32(&inc, 1)
}
