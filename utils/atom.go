package utils

import (
	"strconv"
	"sync/atomic"
)

var uid uint64

func GetUID() string {
	atomic.AddUint64(&uid, 1)
	return strconv.FormatUint(uid, 10)
}
