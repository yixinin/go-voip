package server

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
	"voip/config"
	"voip/protocol"
	"voip/room"

	"go-lib/log"
	"go-lib/registry"
	"go-lib/registry/etcd"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{} // use default options

const (
	ProtocolTCP       = "tcp"
	ProtocolWebSocket = "ws"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Server struct {
	sync.RWMutex
	Rooms           map[int32]*room.Room
	Registry        registry.Registry
	RegistryService *registry.Service
	watcher         registry.Watcher
	tokens          map[string]int64 //[token]uid

	chatClients map[string]protocol.ChatServiceClient

	createRoomChan chan CreateRoom
	closeRoomChan  chan int32
	joinRoomChan   chan JoinRoom
	leaveRoomChan  chan LeaveRoom

	Stop      chan bool
	stopWatch chan bool

	config *config.Config
}

func NewServer(c *config.Config) *Server {
	var regist = etcd.NewRegistry()
	var rooms = make(map[int32]*room.Room, 2)

	var s = &Server{
		Rooms:          rooms,
		config:         c,
		Registry:       regist,
		tokens:         make(map[string]int64, 2*10),
		chatClients:    make(map[string]protocol.ChatServiceClient),
		createRoomChan: make(chan CreateRoom),
		closeRoomChan:  make(chan int32),
		joinRoomChan:   make(chan JoinRoom),
		leaveRoomChan:  make(chan LeaveRoom),

		Stop:      make(chan bool),
		stopWatch: make(chan bool),
	}

	var srv = grpc.NewServer()
	var rs = NewRoomServer(s.createRoomChan, s.joinRoomChan, s.leaveRoomChan, s.closeRoomChan, s.config)
	protocol.RegisterRoomServiceServer(srv, rs)
	var listen, err = net.Listen("tcp", fmt.Sprintf(":%s", s.config.GrpcPort))
	if err != nil {
		log.Error(err)
	}

	go func() {
		err := srv.Serve(listen)
		if err != nil {
			log.Error(err)
		}
	}()
	s.Registry.Init(
		registry.Addrs(s.config.EtcdAddr...),
	)
	watcher, err := s.Registry.Watch()
	if err != nil {
		log.Error(err)
		return s
	}
	s.watcher = watcher
	return s
}

func (s *Server) GetRoom(rid int32) (*room.Room, bool) {
	s.RLock()
	defer s.RUnlock()
	r, ok := s.Rooms[rid]
	return r, ok
}
func (s *Server) AddRoom(r *room.Room) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Rooms[r.RoomId]; ok {
		log.Warnf("repeat room: %d, ignore", r.RoomId)
		return
	}
	s.Rooms[r.RoomId] = r
}

func (s *Server) DelRoom(rid int32) {
	s.Lock()
	defer s.Unlock()
	if r, ok := s.Rooms[rid]; ok {
		close(r.PktChan)
		delete(s.Rooms, rid)
	}
}

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
		case <-s.stopWatch:
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
