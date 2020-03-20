package server

import (
	"context"
	"go-lib/ip"
	"go-lib/log"
	"time"
	"voip/protocol"
)

func (s *Server) LeaveRoom(uid int64) {
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
}

func (s *Server) JoinRoom(uid int64, rid int32, p string) {
	var client protocol.ChatServiceClient
	for _, c := range s.chatClients {
		client = c
	}

	if client != nil {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.JoinRoom(ctx, &protocol.JoinRoomReq{
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
}
