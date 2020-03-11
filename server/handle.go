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
	s.handleReader(reader, writer)
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
	s.handleReader(reader, writer)
}

func (s *Server) handleReader(reader *bufio.Reader, writer *bufio.Writer) {
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
		case <-s.stopTcp:
			return
		case <-s.stopWs:
			return
		default:
			//数据包格式 1+1+2+8 frameType + dataType + dataLength + timeStamp
			var header = make([]byte, 4+8)
			_, err := reader.Read(header)
			length := utils.BytesToUint16(header[2:4])

			var body = make([]byte, length)
			// var buf = make([]byte, 2+8+9600)
			_, err = reader.Read(body)
			if err != nil {
				log.Println("read buffer error:", err)
				return
			}

			var buf = make([]byte, len(header)+len(body))
			copy(buf[:len(header)], header)
			copy(buf[len(header):], body)

			var p = av.NewPacket(buf, uid)
			uidBytes := utils.Int64ToBytes(uid)

			copy(p.Data[4:8+4], uidBytes)
			r, ok := s.Rooms[rid]
			if !ok {
				log.Println("room not exsist", rid)
				return
			}
			// log.Println(len(p.Data))
			r.PktChan <- p
		}
	}
}
