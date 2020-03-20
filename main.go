package main

import (
	"flag"
	"fmt"
	"go-lib/ip"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"voip/config"
	"voip/server"
)

var (
	configPath = flag.String("conf", `config/app.yaml`, "-conf=xxx")
)

func main() {
	var conf, err = config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("load config error: %v", err)
		os.Exit(0)
	}

	var server = server.NewServer(conf)
	server.Serve()

	{ //测试
		createRoom(conf.HttpPort)
		showIP(conf.GrpcPort)
	}

	//监听退出信号
	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c)

	for {
		select {
		case sig := <-c:
			switch sig {
			case os.Interrupt:
				server.Shutdown()
			}
			return
		}
	}
}

func showIP(port string) {
	log.Println("本机IP:", ip.GrpcAddr(port))
}

func createRoom(port string) {
	var body = `{"RoomId":10240,"Users":[{"Uid":102,"VideoPush":true,"AudioPush":true,"Token":"00000000000000000000000000000000"},{"Uid":104,"VideoPush":true,"AudioPush":true,"Token":"00000000000000000000000000000001"}]}`
	http.DefaultClient.Post(fmt.Sprintf("http://localhost:%s/createRoom", port), "application/json", strings.NewReader(body))
}
