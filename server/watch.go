package server

import (
	"go-lib/pool"
	"voip/protocol"

	log "github.com/sirupsen/logrus"

	"go-lib/registry"
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

func (s *Server) Watch(watcher registry.Watcher) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("watcher recovered ", err)
		}
	}()

	services, err := s.Registry.GetService(ChatServiceName)
	if err != nil {
		log.Error(err)
	} else {
		if services != nil {
			for _, srv := range services {
				if srv == nil || srv.Nodes == nil {
					continue
				}
				for _, node := range srv.Nodes {
					if node == nil {
						continue
					}
					var addr = node.Address
					if addr[0] == ':' {
						addr = "chat" + addr
					}
					pool.DefaultGrpcConnPool.AddNode(addr)
					log.Infof("add node %s :%s", srv.Name, addr)
				}
			}
		}

	}

FOR:
	for {
		select {
		case <-s.stop:
			log.Info("exit ...")
			err := s.Registry.Deregister(s.RegistryService)
			if err != nil {
				log.Error(err)
			}
			return
		default:
			res, err := watcher.Next()
			if err != nil || res == nil || res.Service == nil || res.Service.Nodes == nil {
				if err != nil {
					log.Error(err)
				}
				continue FOR
			}
			var name = res.Service.Name
			if name != s.RegistryService.Name {
				switch res.Action {
				case registry.Create.String():
					for _, node := range res.Service.Nodes {
						if node == nil {
							continue FOR
						}
						var addr = node.Address
						if addr[0] == ':' {
							addr = "chat" + addr
						}
						pool.DefaultGrpcConnPool.AddNode(addr)
						log.Infof("----new node %s :%s", name, addr)
					}
				case registry.Update.String():
					for _, node := range res.Service.Nodes {
						if node == nil {
							continue FOR
						}
						var addr = node.Address
						if addr[0] == ':' {
							addr = "chat" + addr
						}
						pool.DefaultGrpcConnPool.AddNode(addr)
						log.Infof("----update node %s :%s", name, addr)
					}
				case registry.Delete.String():
					for _, node := range res.Service.Nodes {
						if node == nil {
							continue FOR
						}
						var addr = node.Address
						if addr[0] == ':' {
							addr = "chat" + addr
						}
						pool.DefaultGrpcConnPool.AddNode(addr)
						log.Infof("----del node %s :%s", name, addr)
					}
				default:
					for _, node := range res.Service.Nodes {
						if node == nil {
							continue FOR
						}
						log.Infof("----not cased, %s node %s :%s", res.Action, name, node.Address)
					}
				}
			}
		}
	}
}
