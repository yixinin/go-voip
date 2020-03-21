package server

import (
	"voip/protocol"
	"voip/room"

	log "github.com/sirupsen/logrus"
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
	Uid    int64
}

func (s *Server) manageRoomUser() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
FOR:
	for {
		select {
		case <-s.stop:
			for _, r := range s.rooms {
				close(r.PktChan)
			}
			return
		case rid := <-s.closeRoomChan:
			s.DelRoom(rid)

		case createRoom := <-s.createRoomChan:

			var r = room.NewRoom(createRoom.RoomId, createRoom.Users)
			s.AddRoom(r)
			s.AddUsers(createRoom.Users, createRoom.RoomId)

		case joinRoom := <-s.joinRoomChan:
			r, ok := s.GetRoom(joinRoom.RoomId)
			if !ok {
				log.Warnf("no such room: %d, ignore", joinRoom.RoomId)
				continue FOR
			}
			r.JoinRoom(joinRoom.User.Uid, nil)
			s.AddUser(joinRoom.User, joinRoom.RoomId)

		case leaveRoom := <-s.leaveRoomChan:
			r, ok := s.GetRoom(leaveRoom.RoomId)
			if !ok {
				log.Warnf("no such room: %d, ignore", leaveRoom.RoomId)
				continue FOR
			}
			r.LeaveRoom(leaveRoom.Uid)
			s.DelUser(leaveRoom.Uid)
		default:
		}
	}
}

func (s *Server) GetRooms() map[int32]*room.Room {
	s.RLock()
	defer s.RUnlock()
	return s.rooms
}

func (s *Server) GetRoom(rid int32) (*room.Room, bool) {
	s.RLock()
	defer s.RUnlock()
	r, ok := s.rooms[rid]
	return r, ok
}
func (s *Server) AddRoom(r *room.Room) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.rooms[r.RoomId]; ok {
		log.Warnf("repeat room: %d, ignore", r.RoomId)
		return
	}
	s.rooms[r.RoomId] = r
}

func (s *Server) DelRoom(rid int32) {
	s.Lock()
	defer s.Unlock()
	if r, ok := s.rooms[rid]; ok {
		close(r.PktChan)
		delete(s.rooms, rid)
	}
}
