package server

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/yixinin/go-voip/config"
	"github.com/yixinin/go-voip/protocol"
	"github.com/yixinin/go-voip/room"
	"github.com/yixinin/go-voip/user"

	log "github.com/sirupsen/logrus"

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
	rooms map[int32]*room.Room
	users map[string]*user.User //[token]uid

	// watcher         registry.Watcher

	chatClients map[string]protocol.ChatServiceClient

	createRoomChan chan CreateRoom
	closeRoomChan  chan int32
	joinRoomChan   chan JoinRoom
	leaveRoomChan  chan LeaveRoom

	stop   chan bool
	stoped bool

	config *config.Config

	udpConn *net.UDPConn
}

func NewServer(c *config.Config) *Server {
	var rooms = make(map[int32]*room.Room, 2)

	var s = &Server{
		rooms:       rooms,
		config:      c,
		users:       make(map[string]*user.User, 2*10),
		chatClients: make(map[string]protocol.ChatServiceClient),

		createRoomChan: make(chan CreateRoom),
		closeRoomChan:  make(chan int32),
		joinRoomChan:   make(chan JoinRoom),
		leaveRoomChan:  make(chan LeaveRoom),

		stop: make(chan bool),
	}

	var srv = grpc.NewServer()
	var rs = NewRoomServer(s.createRoomChan, s.joinRoomChan, s.leaveRoomChan, s.closeRoomChan, s.config)
	protocol.RegisterRoomServiceServer(srv, rs)
	var listen, err = net.Listen("tcp", s.config.GrpcAddr)
	if err != nil {
		log.Error(err)
	}

	go func() {
		err := srv.Serve(listen)
		if err != nil {
			log.Error(err)
		}
	}()
	return s
}

func (s *Server) Shutdown() {
	if s.stoped {
		return
	}
	select {
	case <-s.stop:
		return
	default:
		close(s.stop)
		s.stoped = true
	}
}
