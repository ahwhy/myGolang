package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/week11/socket"
)

//长连接
func main_tcp_client4() {
	ip := "127.0.0.1" //ip换成0.0.0.0和空字符串试试
	port := 5656
	conn, err := net.DialTimeout("tcp4", ip+":"+strconv.Itoa(port), 30*time.Minute)
	socket.CheckError(err)
	fmt.Printf("establish connection to server %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	for { //长连接，即连接建立后进行多轮的读写交互
		request := socket.Request{A: 7, B: 4}
		requestBytes, _ := json.Marshal(request)
		_, err = conn.Write(requestBytes)
		socket.CheckError(err)
		fmt.Printf("write request %s\n", string(requestBytes))
		responseBytes := make([]byte, 256) //初始化后byte数组每个元素都是0
		read_len, err := conn.Read(responseBytes)
		socket.CheckError(err)
		var response socket.Response
		json.Unmarshal(responseBytes[:read_len], &response) //json反序列化时会把0都考虑在内，所以需要指定只读前read_len个字节
		fmt.Printf("receive response: %d\n", response.Sum)
		time.Sleep(1 * time.Second)
	}
}

//先启动tcp_server，再在另外一个终端启动tcp_client
//go run socket/client/tcp_client4.go
