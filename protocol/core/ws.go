package core

import (
	"github.com/gorilla/websocket"
)

type WsConn struct {
	Conn *websocket.Conn
}

func NewWsConn(conn *websocket.Conn) *WsConn {
	return &WsConn{
		Conn: conn,
	}
}

func (c WsConn) Read(buf []byte) (n int, err error) {
	_, r, err := c.Conn.NextReader()
	if err != nil {

		return 0, err
	}
	return r.Read(buf)
}

func (c WsConn) Write(buf []byte) (n int, err error) {
	w, err := c.Conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}
