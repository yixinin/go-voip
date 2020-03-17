package server

import "voip/protocol"

func (s *Server) AddUser(uid, token string) {
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

func (s *Server) DelUser(uid string) {
	s.Lock()
	defer s.Unlock()
	for k, v := range s.tokens {
		if v == uid {
			delete(s.tokens, k)
			break
		}
	}
	if ch, ok := s.stopTcp[uid]; ok {
		ch <- true
		delete(s.stopTcp, uid)
	}
	if ch, ok := s.stopWs[uid]; ok {
		ch <- true
		delete(s.stopWs, uid)
	}
}

func (s *Server) DelToken(token string) {
	s.Lock()
	defer s.Unlock()
	if uid, ok := s.tokens[token]; ok {
		delete(s.tokens, token)
		if ch, ok := s.stopTcp[uid]; ok {
			ch <- true
			delete(s.stopTcp, uid)
		}
		if ch, ok := s.stopWs[uid]; ok {
			ch <- true
			delete(s.stopWs, uid)
		}
	}
	return
}

func (s *Server) GetToken(token string) (uid string, ok bool) {
	s.RLock()
	defer s.RUnlock()
	uid, ok = s.tokens[token]
	return
}
