package rw

import (
	"net"
)

type TcpReaderWriter struct {
	conn   net.Conn
	closed bool
}

func NewTcpReaderWriter(conn net.Conn) *TcpReaderWriter {
	return &TcpReaderWriter{
		conn: conn,
	}
}

func (c *TcpReaderWriter) Read(buf []byte) (n int, err error) {

	return c.conn.Read(buf)
}

func (c *TcpReaderWriter) Write(buf []byte) (n int, err error) {

	return c.conn.Write(buf)
}

func (c *TcpReaderWriter) Close() error {
	if c.closed {
		c.closed = true
		return c.conn.Close()
	}
	return nil
}

func (*TcpReaderWriter) Name() string {
	return ProtocolTCP
}
