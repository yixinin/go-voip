package user

import (
	"net"
	"voip/rw"
)

type RoomUser struct {
	Uid int64
	// RoomId int64
	// Reader bufio.Reader
	Writer    rw.WriterCloser
	Avlible   bool
	VideoPush bool
	AudioPush bool
	Ts        int64
	// Token     string
	UdpAddr *net.UDPAddr // remote udp addr
}

type User struct {
	Token  string
	RoomId int32
	Uid    int64
	Addr   string //chat服务器
}
