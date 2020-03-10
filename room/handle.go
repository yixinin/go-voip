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

FOR:
	for {
		select {
		case <-r.Stop:
			return
		case p := <-r.PktChan:
			if p == nil || len(p.Data) == 0 {
				log.Println("empty data, packet is nil?", p == nil)
				continue FOR
			}
			if p.IsAudio() {
				//音频处理
			} else {
				if p.IsVideo() {
					//视频处理
				}
			}

			// var data = make([]byte, len(p.Data)+8)
			// copy(data, p.Data)
			// uidBytes := utils.Int64ToBytes(p.Uid)
			// copy(data[2:8+2], uidBytes)
			// copy(data[:2], p.Data[:2])

			// // var fromUid = utils.BytesToInt64(data[2 : 8+2])
			// var dataType = data[1]
			// var frameType = data[0]

			// log.Printf("dataType = %d", p.Data[0])

			// log.Println(len(p.Data))
			//推送
			for _, u := range r.Users {
				if u != nil && u.Writer != nil {
					// log.Println(len(p.Data))
					u.Writer.Write(p.Data)
				}
			}
		}
	}
}
