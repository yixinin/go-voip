package room

import (
	"log"
	"voip/utils"
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

	for {
		select {
		case <-r.stop:
			return
		case p := <-r.PktChan:
			if p.IsAudio() {
				//音频处理
			} else {
				if p.IsVideo() {
					//视频处理
				}
			}

			var data = make([]byte, len(p.Data)+8)
			copy(data[2+8:], p.Data[2:])
			uidBytes := utils.Int64ToBytes(p.Uid)
			copy(data[2:8], uidBytes)
			copy(data[:2], p.Data[:2])
			//推送
			for _, u := range r.Users {
				// if u != nil && u.Writer != nil && u.Writer.Alive() {
				// 	u.Writer.Write(&p)
				// }
				u.Writer.Write(data)
			}
		}
	}
}
