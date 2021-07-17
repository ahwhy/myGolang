package main

import (
	"flag"
	"log"

	"github.com/ahwhy/myGolang/week6/homework/simple-http-probe/config"
	"github.com/ahwhy/myGolang/week6/homework/simple-http-probe/http"
)

var (
	configFile string
)

func main() {
	// 传入配置文件路径
	flag.StringVar(&configFile, "c", "simple_http_probe.yml", "config file path")
	// 解析yaml
	conf, err := config.LoadFile(configFile)
	if err != nil {
		log.Printf("[config.Load.error][err:%v]", err)
		return
	}
	log.Printf("配置是：%v", conf)
	// 启动gin
	go http.StartGin(conf)
	select {}
}
