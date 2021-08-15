package main

import (
	"encoding/json"
	"fmt"
	"github.com/ahwhy/myGolang/week11/socket"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main_ws_client() {
	dialer := &websocket.Dialer{}
	header := http.Header{
		"Cookie": []string{"name=zcy"},
	}
	conn, resp, err := dialer.Dial("ws://localhost:5657/add", header) //Dial:握手阶段，会发送一条http请求
	if err != nil {
		fmt.Printf("dial server error:%v\n", err)
		return
	}
	fmt.Println("handshake response header")
	for key, values := range resp.Header {
		fmt.Printf("%s:%s\n", key, values[0])
	}
	// time.Sleep(5 * time.Second)
	defer conn.Close()
	for i := 0; i < 10; i++ {
		request := socket.Request{A: 7, B: 4}
		requestBytes, _ := json.Marshal(request)
		err = conn.WriteJSON(request) //websocket.Conn直接提供发json序列化和反序列化方法
		socket.CheckError(err)
		fmt.Printf("write request %s\n", string(requestBytes))
		var response socket.Response
		err = conn.ReadJSON(&response)
		socket.CheckError(err)
		fmt.Printf("receive response: %d\n", response.Sum)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(30 * time.Second)
}

//go run socket/client/ws_client.go
