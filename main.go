package main

import (
	"flag"
	"log"
	"os"
	"voip/config"
	"voip/server"
)

var (
	configPath = flag.String("conf", `C:\Users\yixin\go\voip\config\app.yaml`, "-conf=xxx")
)

func main() {
	var b byte = 0x1
	// strings.
	log.Printf("%08b & %08b = %08b", b, 0xf, b&0xf)
	log.Printf("1=%08b, 2=%08b", 1, 2)
	var conf, err = config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("load config error: %v", err)
		os.Exit(0)
	}
	var server = server.NewServer(conf)
	server.Serve()
}
