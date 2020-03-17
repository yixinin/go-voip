package server

import (
	"bufio"
	"context"
	"go-lib/ip"
	"go-lib/log"
	"net"
	"strings"
	"time"
	"voip/av"
	"voip/protocol"
	"voip/protocol/core"
	"voip/utils"

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
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	s.handleReader(reader, writer, "tcp")
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
	var wsConn = core.NewWsConn(conn)
	//获取一个连接的reader读取流
	writer := bufio.NewWriter(wsConn)
	reader := bufio.NewReader(wsConn)

	s.handleReader(reader, writer, "ws")
}

func (s *Server) handleHttp(buf []byte) {
	//获取token,roomId
	if len(buf) < 32+4 {
		return
	}
	var token = strings.TrimSpace(string(buf[:32]))
	var rid = utils.BytesToInt32(buf[32:])
	uid, ok := s.GetToken(token)
	if !ok { //鉴权
		log.Warnf("access denied, uid:%s", uid)
		return
	}
	r, ok := s.GetRoom(rid)
	if !ok {
		log.Warnf("access denied,  rid:%d", rid)
		return
	}
	if !r.InRoom(uid) {
		log.Warnf("access denied,rid:%d, uid:%s", rid, uid)
		return
	}
}

func (s *Server) handleReader(reader *bufio.Reader, writer *bufio.Writer, p string) {
	var uid string
	defer func() {
		if r := recover(); r != nil {
			log.Warn("reader serve panic: ", r)
		}
		log.Warn(" Disconnected : ", uid)
		//删除用户
		var client protocol.ChatServiceClient
		for _, c := range s.chatClients {
			client = c
			break
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.LeaveRoom(ctx, &protocol.LeaveRoomReq{
			Uid: uid,
		})
		if err != nil {
			log.Error(err)
		}
	}()

	var header = make([]byte, 2+32+4) //(token + roomid)
	_, err := reader.Read(header)
	if err != nil {
		log.Fatal(err)
		return
	}

	var token = strings.TrimSpace(string(header[2 : 32+2]))
	var rid = utils.BytesToInt32(header[32+2:])

	uid, ok := s.GetToken(token)
	if !ok { //鉴权
		log.Warnf("access denied, uid:%d", uid)
		return
	}
	r, ok := s.GetRoom(rid)
	if !ok {
		log.Warnf("access denied, uid:%d", uid)
		return
	}

	if !r.JoinRoom(uid, writer) {
		log.Warnf("access denied, roomId:%d, uid:%d", rid, uid)
		return
	}

	//成功连接 发送给chatserver
	var client protocol.ChatServiceClient
	for _, c := range s.chatClients {
		client = c
	}

	if client != nil {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.JoinRoom(ctx, &protocol.JoinRoomReq{
			RoomId:   rid,
			Addr:     ip.GrpcAddr(s.config.GrpcPort),
			Protocol: p,
			User: &protocol.RoomUser{
				Uid: uid,
			},
		})
		if err != nil {
			log.Error(err)
		}
	}

	var stop = make(chan bool)
	switch p {
	case "tcp":
		if oldStop, ok := s.stopTcp[uid]; ok {
			oldStop <- true
		}
		s.stopTcp[uid] = stop
	case "ws":
		if oldStop, ok := s.stopWs[uid]; ok {
			oldStop <- true
		}
		s.stopWs[uid] = stop
	}

	//读取数据
	for {
		select {
		case <-stop:
			close(stop)
			return
		default:
			//数据包格式 1+1+4 frameType + dataType + dataLength
			var header = make([]byte, 2+4)
			n1, err := reader.Read(header)
			if err != nil {
				log.Warn("read buffer error:", err)
				return
			}
			length := utils.BytesToUint32(header[2:])

			// log.Warn(header)

			if length == 0 {
				log.Warnf("header length:%d, body expect length:%d", n1, length)
				continue
			}

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

				n2, err := reader.Read(subBody)
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
}
