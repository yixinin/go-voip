package server

import (
	"context"
	"go-lib/ip"
	"time"
	"voip/protocol"

	log "github.com/sirupsen/logrus"
)

func (s *Server) LeaveRoom(uid int64) {
	var client, ok = s.GetRandomRoomClient()
	if !ok {
		return
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
	var client, ok = s.GetRandomRoomClient()
	if !ok {
		return
	}
	if client != nil {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.JoinRoom(ctx, &protocol.JoinRoomReq{
			RoomId: rid,
			Addr:   ip.GetAddr(s.config.GrpcPort),
			User: &protocol.RoomUser{
				Uid: uid,
			},
		})
		if err != nil {
			log.Error(err)
		}
	}
}
