package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yixinin/go-voip/config"
	"github.com/yixinin/go-voip/server"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	configPath = flag.String("conf", `config/app.yaml`, "-conf=xxx")
)

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var conf, err = config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("load config error: %v", err)
		os.Exit(0)
	}

	var server = server.NewServer(conf)
	server.Serve()

	{ //测试
		// createRoom(conf.HttpPort)
		showIP(conf.GrpcHost + conf.GrpcAddr)
	}

	//监听退出信号
	c := make(chan os.Signal, 1)
	//监听所有信号
	signal.Notify(c)
	defer server.Shutdown()
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-c:
			switch sig {
			case os.Interrupt:
				return
			}
			logrus.Info("receive sig", sig)
		}
	}
}

func showIP(port string) {
	// log.Println("本机IP:", ip.GetAddr(port))
}

func createRoom(port string) {
	var body = `{"RoomId":10240,"Users":[{"Uid":102,"VideoPush":true,"AudioPush":true,"Token":"00000000000000000000000000000000"},{"Uid":104,"VideoPush":true,"AudioPush":true,"Token":"00000000000000000000000000000001"}]}`
	http.DefaultClient.Post(fmt.Sprintf("http://localhost:%s/createRoom", port), "application/json", strings.NewReader(body))
}
