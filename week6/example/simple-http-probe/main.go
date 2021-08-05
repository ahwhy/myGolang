package main

import (
	"flag"
	"log"

	"github.com/ahwhy/myGolang/week6/example/simple-http-probe/config"
	"github.com/ahwhy/myGolang/week6/example/simple-http-probe/web"
)

var (
	configFile string
)

func main() {
	// 传入配置文件路径
	flag.StringVar(&configFile, "c", "my_http_probe.yaml", "config file path")
	// 解析yaml
	conf, err := config.LoadFile(configFile)
	if err != nil {
		log.Printf("[config.load.error]")
		return
	}
	log.Printf("导入配置为: %v", conf)

	// 启动gin
	go web.StartGin(conf)
	select {}
}
