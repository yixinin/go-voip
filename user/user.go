package user

import (
	"bufio"
)

type User struct {
	Uid int64
	// RoomId int64
	// Reader bufio.Reader
	Writer    *bufio.Writer
	Avlible   bool
	VideoPush bool
	AudioPush bool
	// Token     string
}
