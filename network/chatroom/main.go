package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ahwhy/myGolang/network/chatroom/protocol"
)

func main() {
	help := flag.Bool("h", false, "this is help")
	port := flag.String("p", "5657", "http service port") // 如果命令行不指定port参数，则默认为5657
	flag.Parse()                                          // 解析命令行输入的port参数
	if *help {
		usage()
		os.Exit(0)
	}

	log.Println("Chatroom is preparing to start...")
	protocol.Registry(port)
}

// usage 使用说明
func usage() {
	fmt.Fprintf(os.Stderr, `chatroom version: 0.0.1
Usage: chatroom [-h] -p <Port>
Options:
`)
	flag.PrintDefaults()
}

/*
go build -o chatroom main.go
./chatroom -p 5657
*/
