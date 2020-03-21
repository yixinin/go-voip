package server

import (
	"go-lib/log"
	"voip/protocol"

	"google.golang.org/grpc"
)

func (s *Server) AddNode(addr string) {
	s.Lock()
	defer s.Unlock()
	var conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return
	}
	client := protocol.NewChatServiceClient(conn)
	s.chatClients[addr] = client
}

func (s *Server) UpdateNode(addr string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.chatClients[addr]; ok {
		return
	}
	var conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return
	}
	client := protocol.NewChatServiceClient(conn)
	s.chatClients[addr] = client
}

func (s *Server) DeleteNode(addr string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.chatClients[addr]; ok {
		delete(s.chatClients, addr)
	}
}

func (s *Server) GetRandomChatClient() (addr string, client protocol.ChatServiceClient) {
	s.RLock()
	defer s.RUnlock()
	for k, v := range s.chatClients {
		return k, v
	}
	return "", nil
}

func (s *Server) GetChatClient(addr string) (client protocol.ChatServiceClient, ok bool) {
	s.RLock()
	defer s.RUnlock()
	client, ok = s.chatClients[addr]
	return
}

func (s *Server) Watch() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	services, err := s.Registry.GetService("live-chat.chat")
	if err == nil {
		for _, srv := range services {
			for _, node := range srv.Nodes {
				s.AddNode(node.Address)
				log.Infof("add node %s :%s", srv.Name, node.Address)
			}
		}
	} else {
		log.Error(err)
	}

	for {
		select {
		case <-s.stop:
			return
		default:

			res, err := s.watcher.Next()
			if err != nil {
				log.Error("netx err:%v", err)
				continue
			}
			var name = res.Service.Name
			if name == "live-chat.chat" {
				switch res.Action {
				case "create":
					for _, node := range res.Service.Nodes {
						s.AddNode(node.Address)
						log.Infof("----new node %s :%s", name, node.Address)
					}
				case "update":
					for _, node := range res.Service.Nodes {
						s.UpdateNode(node.Address)
						log.Infof("----update node %s :%s", name, node.Address)
					}
				case "delete":
					for _, node := range res.Service.Nodes {
						s.DeleteNode(node.Address)
						log.Infof("----del node %s :%s", name, node.Address)
					}
				default:
					for _, node := range res.Service.Nodes {
						log.Infof("----not cased, %s node %s :%s", res.Action, name, node.Address)
					}
				}
			} else {
				for _, node := range res.Service.Nodes {
					log.Infof("%s node %s :%s", res.Action, name, node.Address)
				}
			}
		}
	}
}
