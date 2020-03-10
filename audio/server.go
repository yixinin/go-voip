package audio

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"voip/utils"

	cmap "github.com/orcaman/concurrent-map"
)

type Server struct {
	sync.Mutex
	listener net.Listener
	conns    map[int64]cmap.ConcurrentMap
	msg      map[int64]*Message
	uids     map[int64]*uint64
	online   map[int64]bool
}

type Message struct {
	ch  chan Buffer
	uid int64
}

type Buffer struct {
	buf []byte
}

func NewServer() *Server {
	return &Server{
		conns:  make(map[int64]cmap.ConcurrentMap, 10),
		uids:   make(map[int64]*uint64, 10),
		msg:    make(map[int64]*Message, 2),
		online: make(map[int64]bool, 10),
	}
}

func (s *Server) Setonline(uid int64, online bool) {
	s.Lock()
	defer s.Unlock()
	s.online[uid] = online
}

func (s *Server) GetOnline(uid int64) bool {
	return s.online[uid]
}

func (s *Server) AddUser(uids ...int64) {
	s.Lock()
	defer s.Unlock()
	for _, uid := range uids {
		if _, ok := s.conns[uid]; !ok {
			s.conns[uid] = cmap.New()
			s.msg[uid] = &Message{
				ch:  make(chan Buffer, 1),
				uid: uid,
			}
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
			close(s.msg[uid].ch)
			delete(s.msg, uid)
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
	uid := utils.BytesToInt64(header[:8])
	target := utils.BytesToInt64(header[8:])
	s.AddUser(uid, target)
	s.Setonline(uid, true)

	//发送音频流
	fmt.Printf("uid:%d, target:%d\n", uid, target)
	go s.handleWriter(writer, uid, target)
	//接受音频流
	s.handleReader(reader, uid, target)

}

func (s *Server) handleWriter(writer *bufio.Writer, uid, target int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("rtmp serve panic: ", r)
		}
	}()
	// fmt.Println("handle writer:", target)
	msg := s.msg[target].ch
	for {
		buffer := <-msg
		writer.Write(buffer.buf)
	}
}

func (s *Server) handleReader(reader *bufio.Reader, uid, target int64) {
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
		if s.GetOnline(target) {
			s.msg[uid].ch <- Buffer{buf: buf}
		}

		// var key = atomic.AddUint64(s.uids[target], 1)
		// s.conns[target].Set(strconv.FormatUint(key, 10), buf)

	}
}
