package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/yixinin/go-voip/id"
	"github.com/yixinin/go-voip/rw"

	log "github.com/sirupsen/logrus"
)

const (
	ChatServiceName = "live-chat.chat"
)

func (s *Server) Serve() error {
	for _, v := range s.config.Protocol {
		if v == ProtocolTCP {
			go s.ServeSocket()
		}
		if v == ProtocolWebSocket {
			go s.ServeWs()
		}
	}
	go s.ServeHttp()
	go s.manageRoomUser()
	return nil
}

func (s *Server) ServeSocket() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	//定义一个tcp断点
	var tcpAddr *net.TCPAddr
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", s.config.ListenIp, s.config.TcpPort))
	log.Warn("listen socket/tcp in", tcpAddr.String())
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for {
		netconn, err := tcpListener.Accept()
		if err != nil {
			log.Errorf("accept conn error: %v", err)
			continue
		}
		// conn := core.NewConn(netconn, 4*1024)
		log.Warn("remote:", netconn.RemoteAddr().String(), "local:", netconn.LocalAddr().String())
		go s.handleTcp(netconn)
	}
}

func (s *Server) ServeUDP() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	//定义一个tcp断点
	var udpAddr *net.UDPAddr
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", s.config.ListenIp, s.config.TcpPort))
	if err != nil {
		return
	}
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	s.udpConn = udpListener
	for i := 0; i < 10; i++ {
		go s.handleUdp(udpListener)
	}
}

func (s *Server) ServeWs() {
	var addr = fmt.Sprintf("%s:%s", s.config.ListenIp, s.config.HttpPort)
	http.HandleFunc("/live/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn("upgrade:", err)
			return
		}
		s.handleWs(conn)
	})
	log.Warn("listen ws in", addr+"/live")
}

func (s *Server) ServeHttp() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	http.HandleFunc("/live/http", func(w http.ResponseWriter, r *http.Request) {
		var httpRw = rw.NewHttpReaderWriter(w, r.Body)
		s.handleHttpReader(httpRw)
	})

	http.HandleFunc("/live/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	http.HandleFunc("/live/push", func(w http.ResponseWriter, r *http.Request) {
		s.handleHttp(w, r)
	})
	http.HandleFunc("/createRoom", func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			return
		}

		var createRoom CreateRoom

		err = json.Unmarshal(buf, &createRoom)
		if err != nil {
			log.Error(err)
			return
		}
		if len(createRoom.Users) > 1 {
			if createRoom.RoomId == 0 {
				createRoom.RoomId = id.GenTempID()
			}
			s.createRoomChan <- createRoom
		}

	})

	var addr = fmt.Sprintf("%s:%s", s.config.ListenIp, s.config.HttpPort)
	log.Warn("listen http in", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
