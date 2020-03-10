package server

import (
	"math/rand"
	"time"
	"voip/config"
	"voip/room"

	"github.com/gorilla/websocket"
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
	Rooms  map[int64]*room.Room
	conf   *config.Config
	tokens map[string]int64

	createRoomChan chan CreateRoom
	closeRoomChan  chan int64
	stop           chan bool
}

func NewServer(c *config.Config) *Server {
	return &Server{
		Rooms:          make(map[int64]*room.Room, 2),
		conf:           c,
		tokens:         make(map[string]int64, 2*10),
		createRoomChan: make(chan CreateRoom),
		stop:           make(chan bool),
	}
}
