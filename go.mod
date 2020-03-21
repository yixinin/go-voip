module voip

go 1.14

replace go-lib => ../go-lib

require (
	github.com/go-delve/delve v1.4.0 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.1
	github.com/micro/go-micro/v2 v2.2.0
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6 // indirect
	github.com/sirupsen/logrus v1.4.2
	go-lib v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.2.8
)
