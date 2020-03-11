package room

import (
	"bufio"
	"log"
	"sync"
	"voip/av"
	"voip/user"
)

type Room struct {
	sync.RWMutex
	RoomId  int64
	PktChan chan *av.Packet
	Users   map[int64]*user.User

	// Stop chan bool
}

func NewRoom(id int64, us []*user.User) *Room {
	var room = &Room{
		RoomId: id,
		Users:  make(map[int64]*user.User, len(us)),
		// Stop:    make(chan bool),
		PktChan: make(chan *av.Packet),
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

func (r *Room) Broadcast(p *av.Packet) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Broadcast error:%v", err)
		}
	}()
	r.RLock()
	defer r.RUnlock()

	for _, u := range r.Users {
		if u != nil &&
			u.Uid != p.Uid && //不给自己发
			u.Avlible && //在线
			u.Writer != nil {
			u.Writer.Write(p.Data)
		}
	}
}
