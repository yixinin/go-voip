package server

import (
	"go-lib/log"
	"go-lib/pool"
	"voip/protocol"

	"github.com/micro/go-micro/v2/registry"
)

func (s *Server) GetRoomClient(addr string) (protocol.ChatServiceClient, bool) {
	var conn, ok = pool.DefaultGrpcConnPool.GetConn(addr)
	if !ok {
		return nil, false
	}
	return protocol.NewChatServiceClient(conn), true
}

func (s *Server) GetRandomRoomClient() (protocol.ChatServiceClient, bool) {
	var addr, conn = pool.DefaultGrpcConnPool.GetRandomConn()
	if addr == "" {
		return nil, false
	}
	return protocol.NewChatServiceClient(conn), true
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
				pool.DefaultGrpcConnPool.AddNode(node.Address)
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
			if name != s.RegistryService.Name {
				switch res.Action {
				case registry.Create.String():
					for _, node := range res.Service.Nodes {
						pool.DefaultGrpcConnPool.AddNode(node.Address)
						log.Infof("----new node %s :%s", name, node.Address)
					}
				case registry.Update.String():
					for _, node := range res.Service.Nodes {
						pool.DefaultGrpcConnPool.AddNode(node.Address)
						log.Infof("----update node %s :%s", name, node.Address)
					}
				case registry.Delete.String():
					for _, node := range res.Service.Nodes {
						pool.DefaultGrpcConnPool.AddNode(node.Address)
						log.Infof("----del node %s :%s", name, node.Address)
					}
				default:
					for _, node := range res.Service.Nodes {
						log.Infof("----not cased, %s node %s :%s", res.Action, name, node.Address)
					}
				}
			}
		}
	}
}
