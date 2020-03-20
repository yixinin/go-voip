package server

import "voip/protocol"

func (s *Server) AddUser(uid int64, token string) {
	s.Lock()
	defer s.Unlock()
	s.tokens[token] = uid
}

func (s *Server) AddUsers(users []*protocol.RoomUser) {
	s.Lock()
	defer s.Unlock()
	for _, user := range users {
		s.tokens[user.Token] = user.Uid
	}
}

func (s *Server) DelUser(uid int64) {
	s.Lock()
	defer s.Unlock()
	for k, v := range s.tokens {
		if v == uid {
			delete(s.tokens, k)
			break
		}
	}
}

func (s *Server) DelToken(token string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.tokens[token]; ok {
		delete(s.tokens, token)
	}
	return
}

func (s *Server) GetToken(token string) (uid int64, ok bool) {
	s.RLock()
	defer s.RUnlock()
	uid, ok = s.tokens[token]
	return
}
