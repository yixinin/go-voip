package rw

import (
	"net"
	"sync"
)

type UdpReaderWriter struct {
	sync.Mutex
	buffers    []byte
	BufferChan chan []byte
	udpConn    *net.UDPConn
	addr       *net.UDPAddr
	Closed     bool
}

func NewUdpReaderWriter(addr *net.UDPAddr, conn *net.UDPConn) *UdpReaderWriter {
	var rw = &UdpReaderWriter{
		udpConn:    conn,
		BufferChan: make(chan []byte),
		buffers:    make([]byte, 0),
		addr:       addr,
	}
	return rw
}

func (c *UdpReaderWriter) Read(buf []byte) (n int, err error) {
	c.Lock()
	defer c.Unlock()
	var l = len(buf)
	if l > len(c.buffers) {
		l = len(c.buffers)
		copy(buf, c.buffers)
		c.buffers = make([]byte, 0)
		return l, nil
	}

	copy(buf, c.buffers[:l])
	c.buffers = c.buffers[l:]

	return l, nil
}

func (c *UdpReaderWriter) Write(buf []byte) (n int, err error) {
	if c.addr != nil {
		return c.udpConn.WriteToUDP(buf, c.addr)
	}
	return 0, nil
}

func (c *UdpReaderWriter) Close() error {
	if !c.Closed {
		close(c.BufferChan)
		c.Closed = true
	}
	return nil
}

func (*UdpReaderWriter) Name() string {
	return ProtocolTCP
}

func (w *UdpReaderWriter) Handle() {

	for buf := range w.BufferChan {
		w.Lock()
		w.buffers = append(w.buffers, buf...)
		w.Unlock()
	}
}
