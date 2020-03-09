package audio

import (
	"bufio"
	"fmt"
	"live/utils"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	cmap "github.com/orcaman/concurrent-map"
)

type Server struct {
	sync.Mutex
	listener net.Listener
	conns    map[int64]cmap.ConcurrentMap
	uids     map[int64]*uint64
}

func NewServer() *Server {
	return &Server{
		conns: make(map[int64]cmap.ConcurrentMap, 10),
		uids:  make(map[int64]*uint64, 10),
	}
}

func (s *Server) AddUser(uids ...int64) {
	s.Lock()
	defer s.Unlock()
	for _, uid := range uids {
		if _, ok := s.conns[uid]; !ok {
			s.conns[uid] = cmap.New()
			z := uint64(0)
			s.uids[uid] = &z
			fmt.Println("add user", uid)
		}
	}

}
func (s *Server) DelUser(uids ...int64) {
	s.Lock()
	defer s.Unlock()
	for _, uid := range uids {
		if _, ok := s.conns[uid]; ok {
			delete(s.conns, uid)
			delete(s.uids, uid)
			fmt.Println("del user", uid)
		}
	}
}

func (s *Server) Serve(lis net.Listener) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("rtmp serve panic: ", r)
		}
	}()

	for {
		netconn, err := lis.Accept()
		if err != nil {
			return err
		}
		// conn := core.NewConn(netconn, 4*1024)
		log.Println("remote:", netconn.RemoteAddr().String(), "local:", netconn.LocalAddr().String())
		go s.handleConn(netconn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	var ipStr = conn.RemoteAddr().String()
	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	//获取一个连接的reader读取流
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	var header = make([]byte, 8+8) //(userid + call_userid)
	_, err := reader.Read(header)
	if err != nil {
		log.Fatal(err)
		return
	}

	//读取
	userid := utils.BytesToInt64(header[:8])
	callUserid := utils.BytesToInt64(header[8:])
	s.AddUser(userid, callUserid)

	//发送音频流
	go s.handleWriter(writer, callUserid)
	//接受音频流
	s.handleReader(reader, userid)

}

func (s *Server) handleWriter(writer *bufio.Writer, target int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("rtmp serve panic: ", r)
		}
	}()
	for t := range s.conns[target].IterBuffered() {
		buf := t.Val.([]byte)
		n, err := writer.Write(buf)
		if err != nil {
			log.Println(n, err)
		}
	}
}

func (s *Server) handleReader(reader *bufio.Reader, target int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("rtmp serve panic: ", r)
		}
	}()
	for {
		var buf = make([]byte, 1024*4)
		_, err := reader.Read(buf)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("recieved buf", len(buf))

		var key = atomic.AddUint64(s.uids[target], 1)
		s.conns[target].Set(strconv.FormatUint(key, 10), buf)
	}
}
