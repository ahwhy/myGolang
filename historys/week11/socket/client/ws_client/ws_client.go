package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahwhy/myGolang/historys/week11/socket"
	"github.com/gorilla/websocket"
)

func main() {
	// 建立websocket连接
	dialer := &websocket.Dialer{}
	header := http.Header{ // 构造Request报文头部 Cookie:name=atlantis
		"Cookie": []string{"name=atlantis"},
	}
	conn, resp, err := dialer.Dial("ws://localhost:5657/add", header) // Dial 握手阶段，会发送一条http请求
	if err != nil {
		fmt.Printf("Dial server error:%v\n", err)
		return
	}

	// 打印Response报文头部
	fmt.Println("Handshake response header")
	for key, values := range resp.Header {
		fmt.Printf("%s:%s\n", key, values[0])
	}

	// time.Sleep(5 * time.Second)
	defer conn.Close()
	for i := 0; i < 10; i++ {
		// 发送Request
		request := socket.Request{A: 7, B: 4}
		requestBytes, _ := json.Marshal(request)
		err = conn.WriteJSON(request) //websocket.Conn直接提供发json序列化和反序列化方法
		socket.CheckError(err)
		fmt.Printf("write request %s\n", string(requestBytes))

		// 接收服务端发送Response
		var response socket.Response
		err = conn.ReadJSON(&response)
		socket.CheckError(err)
		fmt.Printf("receive response: %d\n", response.Sum)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(30 * time.Second)
}
