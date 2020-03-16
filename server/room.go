package server

import (
	"go-lib/log"
	"voip/protocol"
	"voip/room"
)

type CreateRoom struct {
	RoomId int32
	Users  []*protocol.RoomUser
}
type JoinRoom struct {
	RoomId int32
	User   *protocol.RoomUser
}

type LeaveRoom struct {
	RoomId int32
	Uid    string
}

func (s *Server) handleCreateRoom() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
FOR:
	for {
		select {
		case <-s.Stop:
			for _, r := range s.Rooms {
				close(r.PktChan)
			}
			for _, stop := range s.stopTcp {
				stop <- true
			}
			for _, stop := range s.stopWs {
				stop <- true
			}
			return
		case rid := <-s.closeRoomChan:
			if r, ok := s.Rooms[rid]; ok {
				close(r.PktChan)
				delete(s.Rooms, rid)
				log.Warnf("closed room, id = %d", rid)
			}

		case createRoom := <-s.createRoomChan:

			if _, ok := s.Rooms[createRoom.RoomId]; ok {
				log.Warnf("repeat room: %d, ignore", createRoom.RoomId)
				continue FOR
			}
			var r = room.NewRoom(createRoom.RoomId, createRoom.Users)
			s.Rooms[r.RoomId] = r
			s.AddUsers(createRoom.Users)

		case joinRoom := <-s.joinRoomChan:
			r, ok := s.Rooms[joinRoom.RoomId]
			if !ok {
				log.Warnf("no such room: %d, ignore", joinRoom.RoomId)
				continue FOR
			}
			r.JoinRoom(joinRoom.User.Uid, nil)
			s.AddUser(joinRoom.User.Uid, joinRoom.User.Token)

		case leaveRoom := <-s.leaveRoomChan:
			r, ok := s.Rooms[leaveRoom.RoomId]
			if !ok {
				log.Warnf("no such room: %d, ignore", leaveRoom.RoomId)
				continue FOR
			}
			r.LeaveRoom(leaveRoom.Uid)
			s.DelUser(leaveRoom.Uid)
		}
	}
}
