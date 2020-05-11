package server

import (
	"log"
	"net"
	"voip/rw"
)

func (s *Server) handleUdp(conn *net.UDPConn) {
	var rws = make(map[string]*rw.UdpReaderWriter, 0)
	for {
		var data = make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
			continue
		}
		v, ok := rws[remoteAddr.String()]
		if !ok {
			v = rw.NewUdpReaderWriter(remoteAddr, conn)
			rws[remoteAddr.String()] = v
			go v.Handle()
			go s.handleReader(v)
		}
		if !v.Closed && v.BufferChan != nil {
			v.BufferChan <- data[:n]
		}
	}
}
