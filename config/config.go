package config

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Protocol []string `yaml:"protocol"`
	TcpPort  string   `yaml:"tcp_port"`
	ListenIp string   `yaml:"listen_ip"`
	HttpPort string   `yaml:"http_port"`
	GrpcAddr string   `yaml:"grpc_addr"`
	GrpcHost string   `yaml:"grpc_host"`
	EtcdAddr []string `yaml:"etcd"`
	Host     string   `yaml:"host"`
}

func ParseConfig(path string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	return &conf, err
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// log.SetReportCaller(true)
}
