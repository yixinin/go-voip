package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"voip/av"
	"voip/protocol/core"
	"voip/utils"

	"github.com/gorilla/websocket"
)

const (
	READSIZE = 4096
)

func (s *Server) handleRpc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("pls use POST"))
		return
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("read message error: %s", err)))
		return
	}
	var cr CreateRoom
	err = json.Unmarshal(buf, &cr)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unmarshal message error: %s", err)))
		return
	}
	s.createRoomChan <- cr
	w.Write([]byte("create room success"))
	return
}

func (s *Server) handleTcp(conn net.Conn) {
	var ipStr = conn.RemoteAddr().String()
	defer func() {
		if err := recover(); err != nil {
			log.Println("tcp exited, err:", err)
		}
		log.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	//获取一个连接的reader读取流
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	//接受音频流
	var stop = make(chan bool)
	s.stopTcp = append(s.stopTcp, stop)
	s.handleReader(reader, writer, stop)
}

func (s *Server) handleWs(conn *websocket.Conn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		if err := recover(); err != nil {
			log.Panicln("ws exited err:", err)
		}
		log.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()
	var wsConn = core.NewWsConn(conn)
	//获取一个连接的reader读取流
	writer := bufio.NewWriter(wsConn)
	reader := bufio.NewReader(wsConn)

	//接受音频流
	var stop = make(chan bool)
	s.stopWs = append(s.stopWs, stop)
	s.handleReader(reader, writer, stop)
}

func (s *Server) handleReader(reader *bufio.Reader, writer *bufio.Writer, stop chan bool) {
	var uid int64
	defer func() {
		if r := recover(); r != nil {
			log.Println("reader serve panic: ", r)
		}
		log.Println(" Disconnected : ", uid)
	}()

	var header = make([]byte, 2+32+8) //(userid + roomid)
	_, err := reader.Read(header)
	if err != nil {
		log.Fatal(err)
		return
	}

	var token = strings.TrimSpace(string(header[2 : 32+2]))
	var rid = utils.BytesToInt64(header[32+2:])

	uid, ok := s.tokens[token]
	if !ok { //鉴权
		log.Printf("access denied, uid:%d", uid)
		return
	}

	if !s.Rooms[rid].JoinRoom(uid, writer) {
		log.Printf("access denied, roomId:%d, uid:%d", rid, uid)
		return
	}

	//读取数据
	for {
		select {
		case <-stop:
			return
		default:
			//数据包格式 1+1+4 frameType + dataType + dataLength
			var header = make([]byte, 2+4)
			n1, err := reader.Read(header)
			if err != nil {
				log.Println("read buffer error:", err)
				return
			}
			length := utils.BytesToUint32(header[2:])

			log.Println(header)

			if length == 0 {
				log.Printf("header length:%d, body expect length:%d", n1, length)
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
					log.Printf("body length:%d, body expect length:%d, body read length:%d", n1, length, n2)
					return
				}

				copy(body[read:read+n2], subBody)
				read += n2
			}

			var buf = make([]byte, len(header)+len(body))
			copy(buf[:len(header)], header) //复制头
			copy(buf[len(header):], body)   //复制body

			var p = av.NewPacket(buf, uid)

			r, ok := s.Rooms[rid]
			if !ok {
				log.Println("room not exsist", rid)
				return
			}

			r.PktChan <- p
		}
	}
}
