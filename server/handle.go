package server

import (
	"go-lib/utils"
	"net"
	"net/http"
	"voip/av"
	"voip/rw"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	READSIZE = 4096
)

var (
	LoginPacket = 0
)

func (s *Server) handleTcp(conn net.Conn) {
	var ipStr = conn.RemoteAddr().String()
	defer func() {
		if err := recover(); err != nil {
			log.Warn("tcp exited, err:", err)
		}
		log.Warn(" Disconnected : " + ipStr)
		conn.Close()
	}()

	//获取一个连接的reader读取流
	// writer := bufio.NewWriter(conn)
	// reader := bufio.NewReader(conn)
	var tcpRw = rw.NewTcpReaderWriter(conn)

	s.handleReader(tcpRw)
}

func (s *Server) handleWs(conn *websocket.Conn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("ws exited err:", err)
		}
		log.Warn(" Disconnected : " + ipStr)
		conn.Close()
	}()
	var wsRw = rw.NewWsReaderWriter(conn)
	//获取一个连接的reader读取流
	// writer := bufio.NewWriter(wsConn)
	// reader := bufio.NewReader(wsConn)

	s.handleReader(wsRw)
}

func (s *Server) HandleHttp(w http.ResponseWriter, r *http.Request) {
	var httpRw = rw.NewHttpReaderWriter(w, r.Response.Body)
	s.handleHttpReader(httpRw)
}

func (s *Server) handleHttpReader(readerWriter rw.ReaderWriterCloser) {
	uid, rid, ok := s.Auth(readerWriter)
	if !ok {
		return
	}
	var tsBuf = make([]byte, 8)
	readerWriter.Read(tsBuf)
	var ts = utils.BytesToUint64(tsBuf)
	r, ok := s.GetRoom(rid)
	if ok {
		gop := r.GetGopCache(uid, ts)
		if gop == nil {
			return
		}
		var vl = len(gop.Video)
		var al = len(gop.Audio)

		for i := 0; i < al || i < vl; i++ {
			if i < al {
				readerWriter.Write(gop.Audio[i].Data)
			}
			if i < vl {
				readerWriter.Write(gop.Video[i].Data)
			}
		}

	}
}

func (s *Server) handleReader(readerWriter rw.ReaderWriterCloser) {
	var (
		uid int64
		rid int32
		ok  bool
	)

	defer func() {
		if r := recover(); r != nil {
			log.Error("reader serve panic: ", r)
		}
		log.Warn(" Disconnected : ", uid)
		//删除用户
		s.LeaveRoom(uid)
	}()

	uid, rid, ok = s.Auth(readerWriter)
	if !ok {
		return
	}

	//读取数据
	s.hanldePacket(readerWriter, rid, uid)
}

func (s *Server) hanldePacket(readerWriter rw.ReaderCloser, rid int32, uid int64) {
	for {
		//数据包格式 1+1+4 frameType + dataType + dataLength + timestamp
		var header = make([]byte, 2+4+8)
		n1, err := readerWriter.Read(header)
		if err != nil {
			log.Warn("read buffer error:", err)
			return
		}
		length := utils.BytesToUint32(header[2 : 2+4])

		// log.Warn(header)

		if length == 0 {
			log.Warnf("header length:%d, body expect length:%d", n1, length)
			continue
		}
		// ts := utils.BytesToUint64(header[2+4:])
		var body = make([]byte, length)

		read := 0

		for read < int(length) {
			var currentLength int
			if unRead := int(length) - read; unRead > READSIZE {
				currentLength = READSIZE
			} else {
				currentLength = unRead
			}
			var subBody = make([]byte, currentLength)

			n2, err := readerWriter.Read(subBody)
			if err != nil {
				log.Warnf("body length:%d, body expect length:%d, body read length:%d", n1, length, n2)
				return
			}

			copy(body[read:read+n2], subBody)
			read += n2
		}

		var buf = make([]byte, len(header)+len(body))
		copy(buf[:len(header)], header) //复制头
		copy(buf[len(header):], body)   //复制body

		var p = av.NewPacket(buf, uid)

		r, ok := s.GetRoom(rid)
		if !ok {
			log.Warn("room not exsist", rid)
			return
		}

		r.PktChan <- p
	}
}
