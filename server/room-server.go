package server

import (
	"context"
	"fmt"

	"github.com/yixinin/go-voip/config"
	"github.com/yixinin/go-voip/id"
	"github.com/yixinin/go-voip/protocol"
)

type RoomServer struct {
	createRoomCh chan CreateRoom
	joinRoomCh   chan JoinRoom
	leaveRoomCh  chan LeaveRoom
	closeRoomCh  chan int32
	config       *config.Config
}

func NewRoomServer(create chan CreateRoom, join chan JoinRoom, leave chan LeaveRoom, clos chan int32, c *config.Config) *RoomServer {
	var s = &RoomServer{
		createRoomCh: create,
		joinRoomCh:   join,
		leaveRoomCh:  leave,
		closeRoomCh:  clos,
		config:       c,
	}

	return s
}

func (s *RoomServer) CreateRoom(ctx context.Context, req *protocol.CreateRoomReq) (ack *protocol.CreateRoomAck, err error) {
	ack = new(protocol.CreateRoomAck)
	ack.Header = &protocol.CallAckHeader{
		Code: 200,
		Msg:  "Success",
	}
	var rid = id.GenTempID()
	s.createRoomCh <- CreateRoom{
		Users:  req.Users,
		RoomId: rid,
	}

	ack.RoomId = rid
	ack.TcpAddr = fmt.Sprintf("%s:%s", s.config.Host, s.config.TcpPort)
	ack.WsAddr = fmt.Sprintf("ws://%s:%s/ws/live", s.config.Host, s.config.HttpPort)
	ack.HttpAddr = fmt.Sprintf("http://%s:%s/http/live", s.config.Host, s.config.HttpPort)
	return
}

func (s *RoomServer) JoinRoom(ctx context.Context, req *protocol.JoinRoomReq) (ack *protocol.JoinRoomAck, err error) {
	ack = new(protocol.JoinRoomAck)
	ack.Header = &protocol.CallAckHeader{
		Code: 200,
		Msg:  "Success",
	}
	s.joinRoomCh <- JoinRoom{
		RoomId: req.RoomId,
		User:   req.User,
	}
	return
}

func (s *RoomServer) LeaveRoom(ctx context.Context, req *protocol.LeaveRoomReq) (ack *protocol.LeaveRoomAck, err error) {
	ack = new(protocol.LeaveRoomAck)
	ack.Header = &protocol.CallAckHeader{
		Code: 200,
		Msg:  "Success",
	}
	s.leaveRoomCh <- LeaveRoom{
		RoomId: req.RoomId,
		Uid:    req.Uid,
	}
	return
}

func (s *RoomServer) DiscardRoom(ctx context.Context, req *protocol.DiscardRoomReq) (ack *protocol.DiscardRoomAck, err error) {
	ack = new(protocol.DiscardRoomAck)
	ack.Header = &protocol.CallAckHeader{
		Code: 200,
		Msg:  "Success",
	}
	s.closeRoomCh <- req.RoomId
	return
}
