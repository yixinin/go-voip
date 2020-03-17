package room

import (
	"bufio"
	"log"
	"sync"
	"voip/av"
	"voip/protocol"
	"voip/user"
)

type Room struct {
	sync.RWMutex
	RoomId  int32
	PktChan chan *av.Packet
	Users   map[string]*user.User

	// Stop chan bool
}

func NewRoom(id int32, us []*protocol.RoomUser) *Room {
	var room = &Room{
		RoomId: id,
		Users:  make(map[string]*user.User, len(us)),
		// Stop:    make(chan bool),
		PktChan: make(chan *av.Packet, 100),
	}
	for _, u := range us {
		room.Users[u.Uid] = &user.User{
			Uid:       u.Uid,
			VideoPush: u.VideoPush,
			AudioPush: u.AudioPush,
		}
	}
	go room.handlePacket()
	return room
}

func (r *Room) JoinRoom(uid string, writer *bufio.Writer) bool {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[uid]; !ok {
		return false
	}
	r.Users[uid].Writer = writer
	r.Users[uid].Avlible = true
	return true
}

func (r *Room) InRoom(uid string) bool {
	_, ok := r.Users[uid]
	return ok
}

func (r *Room) LeaveRoom(uid string) {
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

	log.Printf("header:%v, uid:%d", p.Data[:6], p.Uid)

	for _, u := range r.Users {
		if u != nil &&
			u.Uid != p.Uid && //不给自己发
			u.Avlible && //在线
			u.Writer != nil {
			u.Writer.Write(p.Data)
		}
	}
}
