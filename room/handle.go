package room

import (
	"log"
)

// func (r *Room) handleUserAction() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Fatalf("handleUserAction recovered, roomId:%d, err:%v", r.RoomId, err)
// 		}
// 	}()

// 	for {
// 		select {
// 		case <-r.stop:
// 			return
// 		case ua := <-r.userChan:
// 			if ua.IsJoin {
// 				r.JoinRoom(&ua.User)
// 			} else {
// 				r.LeaveRoom(ua.User.Id)
// 			}
// 		}
// 	}
// }

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
