package server

import (
	"context"
	"fmt"
	"go-lib/ip"
	"voip/config"
	"voip/protocol"
	"voip/utils"
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
	var rid = utils.GetRoomID()
	s.createRoomCh <- CreateRoom{
		Users:  req.Users,
		RoomId: utils.GetRoomID(),
	}
	ack.RoomId = rid
	ack.TcpAddr = ip.GetAddr(s.config.TcpPort)
	ack.WsAddr = fmt.Sprintf("ws://%s/ws/live", ip.GetAddr(s.config.HttpPort))
	ack.HttpAddr = fmt.Sprintf("http://%s/http/live", ip.GetAddr(s.config.HttpPort))
	return
}

func (s *RoomServer) JoinRoom(ctx context.Context, req *protocol.JoinRoomReq) (ack *protocol.JoinRoomAck, err error) {
	s.joinRoomCh <- JoinRoom{
		RoomId: req.RoomId,
		User:   req.User,
	}
	return
}

func (s *RoomServer) LeaveRoom(ctx context.Context, req *protocol.LeaveRoomReq) (ack *protocol.LeaveRoomAck, err error) {
	s.leaveRoomCh <- LeaveRoom{
		RoomId: req.RoomId,
		Uid:    req.Uid,
	}
	return
}

func (s *RoomServer) DiscardRoom(ctx context.Context, req *protocol.DiscardRoomReq) (ack *protocol.DiscardRoomAck, err error) {
	s.closeRoomCh <- req.RoomId
	return
}
