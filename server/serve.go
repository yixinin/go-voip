package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

const ()

func (s *Server) Serve() error {

	for _, v := range s.conf.Protocol {
		if v == ProtocolTCP {
			go s.ServeSocket()
		}
		if v == ProtocolTCP {
			s.ServeWs()
		}
	}
	go s.ServeHttp()
	go s.handleCreateRoom()
	return nil
}

func (s *Server) ServeSocket() {
	//定义一个tcp断点
	var tcpAddr *net.TCPAddr
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", s.conf.ListenIp, s.conf.TcpPort))
	log.Println("listen socket/tcp in", tcpAddr.String())
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for {
		netconn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("accept conn error: %v", err)
			continue
		}
		// conn := core.NewConn(netconn, 4*1024)
		log.Println("remote:", netconn.RemoteAddr().String(), "local:", netconn.LocalAddr().String())
		go s.handleTcp(netconn)
	}
}

func (s *Server) ServeWs() {
	var addr = fmt.Sprintf("%s:%s", s.conf.ListenIp, s.conf.HttpPort)
	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		s.handleWs(conn)
	})
	log.Println("listen ws in", addr+"/live")
}

func (s *Server) ServeHttp() {
	var addr = fmt.Sprintf("%s:%s", s.conf.ListenIp, s.conf.HttpPort)
	http.HandleFunc("/createRoom", s.handleRpc)

	log.Println("listen http in", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
