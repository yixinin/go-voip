package server

import (
	"voip/protocol"
	"voip/user"

	log "github.com/sirupsen/logrus"
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
	log.Warnf("add user, token=%s, uid = %d", u.Token, u.Uid)
}

func (s *Server) AddUsers(users []*protocol.RoomUser, rid int32) {
	s.Lock()
	defer s.Unlock()
	for _, u := range users {
		log.Warnf("add user, token=%s, uid = %d", u.Token, u.Uid)
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
			log.Warnf("del user, token=%s, uid = %d", v.Token, uid)
			delete(s.users, k)
			break
		}
	}
}

func (s *Server) DelToken(token string) {
	s.Lock()
	defer s.Unlock()
	if u, ok := s.users[token]; ok {
		log.Warnf("del user, token=%s, uid = %d", token, u.Uid)
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
