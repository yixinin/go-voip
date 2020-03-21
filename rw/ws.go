package rw

import (
	"github.com/gorilla/websocket"
)

type WsReaderWriter struct {
	conn   *websocket.Conn
	closed bool
}

func NewWsReaderWriter(conn *websocket.Conn) *WsReaderWriter {
	return &WsReaderWriter{
		conn: conn,
	}
}

func (c *WsReaderWriter) Read(buf []byte) (n int, err error) {
	_, r, err := c.conn.NextReader()
	if err != nil {

		return 0, err
	}
	return r.Read(buf)
}

func (c *WsReaderWriter) Write(buf []byte) (n int, err error) {
	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

func (c *WsReaderWriter) Close() error {
	if c.closed {
		c.closed = true
		return c.conn.Close()
	}
	return nil
}

func (*WsReaderWriter) Name() string {
	return ProtocolWS
}
