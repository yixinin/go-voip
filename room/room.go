package room

import (
	"bufio"
	"sync"
	"voip/av"
	"voip/user"
)

type Room struct {
	sync.RWMutex
	RoomId  int64
	PktChan chan *av.Packet
	Users   map[int64]*user.User

	Stop chan bool
}

func NewRoom(id int64, us []*user.User) *Room {
	var room = &Room{
		RoomId:  id,
		Users:   make(map[int64]*user.User, len(us)),
		Stop:    make(chan bool),
		PktChan: make(chan *av.Packet, 10),
	}
	for _, u := range us {
		room.Users[u.Uid] = u
	}
	go room.handlePacket()
	return room
}

func (r *Room) JoinRoom(uid int64, writer *bufio.Writer) bool {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[uid]; !ok {
		return false
	}
	r.Users[uid].Writer = writer
	r.Users[uid].Avlible = true
	return true
}

func (r *Room) LeaveRoom(uid int64) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[uid]; ok {
		r.Users[uid].Writer = nil
		r.Users[uid].Avlible = false
	}
}
