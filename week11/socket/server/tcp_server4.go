package main

import (
	"encoding/json"
	"fmt"
	"github.com/ahwhy/myGolang/week11/socket"
	"net"
	"strconv"
	"time"
)

func handleRequest2(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) //30秒后conn.Read会报出i/o timeout
	defer conn.Close()
	for { //长连接，即连接建立后进行多轮的读写交互
		requestBytes := make([]byte, 256) //初始化后byte数组每个元素都是0
		read_len, err := conn.Read(requestBytes)
		if err != nil {
			fmt.Printf("read from socket error: %s\n", err.Error())
			break //到达deadline后，退出for循环，关闭连接。client再用这个连接读写会发生错误
		}
		fmt.Printf("receive request %s\n", string(requestBytes)) //[]byte转string时，0后面的会自动被截掉

		var request socket.Request
		json.Unmarshal(requestBytes[:read_len], &request) //json反序列化时会把0都考虑在内，所以需要指定只读前read_len个字节
		response := socket.Response{Sum: request.A + request.B}

		responseBytes, _ := json.Marshal(response)
		_, err = conn.Write(responseBytes)
		socket.CheckError(err)
		fmt.Printf("write response %s\n", string(responseBytes))
	}
}

//长连接
func main_tcp_server4() {
	ip := "127.0.0.1" //ip换成0.0.0.0和空字符串试试
	port := 5656
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(port))
	socket.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	socket.CheckError(err)
	fmt.Println("waiting for client connection ......")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("establish connection to client %s\n", conn.RemoteAddr().String()) //操作系统会随机给客户端分配一个49152~65535上的端口号
		go handleRequest2(conn)
	}
}

//go run socket/server/tcp_server4.go
