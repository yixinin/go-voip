package room

import (
	"log"
	"sync"
	"voip/av"
	"voip/cache"
	"voip/protocol"
	"voip/rw"
	"voip/user"
)

type Room struct {
	sync.RWMutex
	RoomId  int32
	PktChan chan *av.Packet
	cache   *cache.Cache
	Users   map[int64]*user.RoomUser

	// Stop chan bool
}

func NewRoom(id int32, us []*protocol.RoomUser) *Room {
	var room = &Room{
		RoomId: id,
		Users:  make(map[int64]*user.RoomUser, len(us)),
		// Stop:    make(chan bool),
		PktChan: make(chan *av.Packet, 100),
		cache:   cache.NewCache(),
	}
	for _, u := range us {
		room.Users[u.Uid] = &user.RoomUser{
			Uid:       u.Uid,
			VideoPush: u.VideoPush,
			AudioPush: u.AudioPush,
		}
	}
	go room.handlePacket()
	return room
}

func (r *Room) JoinRoom(uid int64, readerWriter rw.ReaderWriterCloser) bool {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[uid]; !ok {
		return false
	}
	r.Users[uid].Writer = readerWriter
	r.Users[uid].Avlible = true
	return true
}

func (r *Room) InRoom(uid int64) bool {
	_, ok := r.Users[uid]
	return ok
}

func (r *Room) LeaveRoom(uid int64) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[uid]; ok {
		r.Users[uid].Writer.Close()
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
	//缓存
	r.cache.Put(p)
}
