package user

import (
	"voip/rw"
)

type User struct {
	Uid int64
	// RoomId int64
	// Reader bufio.Reader
	Writer    rw.WriterCloser
	Avlible   bool
	VideoPush bool
	AudioPush bool
	// Token     string
}
