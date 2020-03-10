package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func (s *Server) Serve() error {

	switch s.conf.Protocol {
	case "tcp":
		go s.ServeSocket()
	default:
		log.Fatalln("unsurrport protocol")
		return fmt.Errorf("unsurrport protocol")
	}
	go s.ServeHttp()
	s.handleCreateRoom()
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

func (s *Server) ServeHttp() {
	var addr = fmt.Sprintf("%s:%s", s.conf.ListenIp, s.conf.HttpPort)
	http.HandleFunc("/createRoom", s.handleRpc)
	if s.conf.Protocol == "ws" {
		http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Print("upgrade:", err)
				return
			}
			s.handleWs(conn)
		})

	}
	log.Println("listen http/ws in", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
