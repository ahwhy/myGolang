package main

import (
	"flag"
	"log"

	"github.com/ahwhy/myGolang/http_probe/config"
	"github.com/ahwhy/myGolang/http_probe/http"
)

var (
	configFile string
)

func main() {
	// 传入配置文件路径
	flag.StringVar(&configFile, "c", "./etc/http_probe.yaml", "config file path")
	// 解析yaml
	conf, err := config.LoadFile(configFile)
	if err != nil {
		log.Printf("[config.load.error]")
		return
	}
	log.Printf("导入配置为: %v", conf)

	// 启动gin
	go http.StartGin(conf)
	select {}
}
