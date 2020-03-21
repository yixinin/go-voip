package server

import (
	"voip/protocol"
	"voip/user"
)

func (s *Server) AddUser(u *protocol.RoomUser, rid int32) {
	s.Lock()
	defer s.Unlock()
	s.users[u.Token] = &user.User{
		RoomId: rid,
		Token:  u.Token,
		Uid:    u.Uid,
		Addr:   u.Addr,
	}
}

func (s *Server) AddUsers(users []*protocol.RoomUser, rid int32) {
	s.Lock()
	defer s.Unlock()
	for _, u := range users {
		s.users[u.Token] = &user.User{
			Uid:    u.Uid,
			RoomId: rid,
			Token:  u.Token,
			Addr:   u.Addr,
		}
	}
}

func (s *Server) DelUser(uid int64) {
	s.Lock()
	defer s.Unlock()
	for k, v := range s.users {
		if v.Uid == uid {
			delete(s.users, k)
			break
		}
	}
}

func (s *Server) DelToken(token string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.users[token]; ok {
		delete(s.users, token)
	}
	return
}

func (s *Server) GetUser(token string) (u *user.User, ok bool) {
	s.RLock()
	defer s.RUnlock()
	u, ok = s.users[token]
	return
}

func (s *Server) GetUsers() map[string]*user.User {
	s.RLock()
	defer s.RUnlock()
	return s.users
}
