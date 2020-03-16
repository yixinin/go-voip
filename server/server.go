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
	Rooms    map[int32]*room.Room
	Registry registry.Registry
	watcher  registry.Watcher
	conf     *config.Config
	tokens   map[string]string //[token]uid

	chatClients map[string]protocol.ChatServiceClient

	createRoomChan chan CreateRoom
	closeRoomChan  chan int32
	joinRoomChan   chan JoinRoom
	leaveRoomChan  chan LeaveRoom

	Stop      chan bool
	stopWatch chan bool
	stopTcp   []chan bool
	stopWs    []chan bool
	config    *config.Config
}

func NewServer(c *config.Config) *Server {
	var regist = etcd.NewRegistry()
	var s = &Server{
		Rooms:          make(map[int32]*room.Room, 2),
		conf:           c,
		Registry:       regist,
		tokens:         make(map[string]string, 2*10),
		chatClients:    make(map[string]protocol.ChatServiceClient),
		createRoomChan: make(chan CreateRoom),
		closeRoomChan:  make(chan int32),
		joinRoomChan:   make(chan JoinRoom),
		leaveRoomChan:  make(chan LeaveRoom),

		Stop:      make(chan bool),
		stopTcp:   make([]chan bool, 0, 100),
		stopWs:    make([]chan bool, 0, 100),
		stopWatch: make(chan bool),
	}

	var srv = grpc.NewServer()
	var rs = NewRoomServer(s.createRoomChan, s.joinRoomChan, s.leaveRoomChan, s.closeRoomChan)
	protocol.RegisterRoomServiceServer(srv, rs)
	var listen, err = net.Listen("tcp", fmt.Sprintf(":%s", s.config.GrpcPort))
	if err != nil {
		log.Error(err)
	}
	srv.Serve(listen)
	watcher, err := regist.Watch()
	if err != nil {
		log.Error(err)
		return s
	}
	s.watcher = watcher
	go s.Watch()
	return s
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

	for {
		select {
		case <-s.stopWatch:
			return
		default:
			res, err := s.watcher.Next()
			if err != nil {
				log.Error(err)
				continue
			}
			var name = res.Service.Name
			if name == "live-chat.chat" {
				switch res.Action {
				case "create":
					for _, node := range res.Service.Nodes {
						s.AddNode(node.Address)
					}
				case "update":
					for _, node := range res.Service.Nodes {
						s.UpdateNode(node.Address)
					}
				case "delete":
					for _, node := range res.Service.Nodes {
						s.DeleteNode(node.Address)
					}
				}

			}
		}
	}
}
