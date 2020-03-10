package server

import (
	"log"
	"voip/room"
	"voip/user"
)

type UserInfo struct {
	Uid       int64
	VideoPush bool
	AudioPush bool
	Token     string
}

type CreateRoom struct {
	RoomId int64
	Users  []*UserInfo
}

func (s *Server) handleCreateRoom() {
FOR:
	for {
		select {
		case <-s.Stop:
			for _, r := range s.Rooms {
				r.Stop <- true
			}
			return
		case rid := <-s.closeRoomChan:
			if _, ok := s.Rooms[rid]; ok {
				delete(s.Rooms, rid)
				log.Printf("closed room, id = %d", rid)
			}

		case createRoom := <-s.createRoomChan:
			if _, ok := s.Rooms[createRoom.RoomId]; ok {
				log.Printf("repeat room: %d, ignore", createRoom.RoomId)
				continue FOR
			}
			var uids = make([]int64, 0, len(createRoom.Users))
			var us = make([]*user.User, 0, len(createRoom.Users))
			for _, userInfo := range createRoom.Users {
				var u = &user.User{
					Uid:       userInfo.Uid,
					VideoPush: userInfo.VideoPush,
					AudioPush: userInfo.AudioPush,
				}
				us = append(us, u)
				uids = append(uids, u.Uid)
				if _, ok := s.tokens[userInfo.Token]; !ok {
					s.tokens[userInfo.Token] = u.Uid
				}
			}

			s.Rooms[createRoom.RoomId] = room.NewRoom(createRoom.RoomId, us)
			log.Printf("created room, id = %d, users : %v", createRoom.RoomId, uids)
		}
	}
}
