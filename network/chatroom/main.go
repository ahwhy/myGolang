package main

import (
	"flag"

	"github.com/ahwhy/myGolang/network/chatroom/protocol"
)

func main() {
	port := flag.String("port", "5657", "http service port") //如果命令行不指定port参数，则默认为5657
	flag.Parse()                                             //解析命令行输入的port参数

	protocol.Run(port)
}

//go run chat_room/*.go --port 5657
