package room

import (
	"go-lib/log"
)

func (r *Room) handlePacket() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("handlePacket recovered, roomId:%d, err:%v", r.RoomId, err)
		}
	}()

	for p := range r.PktChan {
		//推送
		go r.Broadcast(p)
	}
	return
}
