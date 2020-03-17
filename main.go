package main

import (
	"flag"
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
	showIP(conf.GrpcPort)
	//监听退出信号
	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c)

	for {
		select {
		case sig := <-c:
			switch sig {
			case os.Interrupt:
				server.Stop <- true
			}
			return
		}
	}
}

func showIP(port string) {
	log.Println("本机IP:", ip.GrpcAddr(port))
}

func createRoom() {
	var body = `{"RoomId":10240,"Users":[{"Uid":102,"VideoPush":true,"AudioPush":true,"Token":"00000000000000000000000000000000"},{"Uid":104,"VideoPush":true,"AudioPush":true,"Token":"00000000000000000000000000000001"}]}`
	http.DefaultClient.Post("http://localhost:9902/createRoom", "application/json", strings.NewReader(body))
}
