package main

import (
	"live/audio"
	"log"
	"net"
)

//tcp server 服务端代码

func main() {
	//定义一个tcp断点
	var tcpAddr *net.TCPAddr
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "0.0.0.0:9999")
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	var server = audio.NewServer()
	err := server.Serve(tcpListener)
	if err != nil {
		log.Fatal("serve error", err)
	}
}
