package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

/*
 - 需要一个解析yaml的配置
	- http的listen的地址
	- 探测超时时间
*/

var GlobalTwSec int

type Config struct {
	HttpListenAddr         string `yaml:"http_listen_addr"`
	HttpProbeTimeoutSecond int    `yaml:"http_probe_timeout_second"`
}

func Load(in []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(in, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := Load(content)
	if err != nil {
		log.Printf("load.yaml.errorr: %v", err)
	}
	if cfg.HttpProbeTimeoutSecond == 0 {
		GlobalTwSec = 5
	} else {
		GlobalTwSec = cfg.HttpProbeTimeoutSecond
	}
	return cfg, nil
}
