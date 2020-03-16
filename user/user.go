package user

import (
	"bufio"
)

type User struct {
	Uid string
	// RoomId int64
	// Reader bufio.Reader
	Writer    *bufio.Writer
	Avlible   bool
	VideoPush bool
	AudioPush bool
	// Token     string
}
