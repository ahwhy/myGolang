package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/week11/socket"
)

// 处理Request
func handleLongRequest(conn net.Conn) {
	// 设置超时时间，30秒后 conn.Read 会报出 i/o timeout
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	
	// 中断TCP连接
	defer conn.Close()

	//长连接，即连接建立后进行多轮的读写交互
	for {
		// 获取客户端Request
		requestBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
		read_len, err := conn.Read(requestBytes)
		if err != nil {
			fmt.Printf("Read from socket error: %s\n", err.Error())
			break // 到达deadline后，退出for循环，关闭连接；client再用这个连接读写会发生错误
		}
		fmt.Printf("Receive request %s\n", string(requestBytes)) // []byte转string时，0后面的会自动被截掉

		// 返回Response
		var request socket.Request
		json.Unmarshal(requestBytes[:read_len], &request) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节
		response := socket.Response{Sum: request.A + request.B}
		responseBytes, _ := json.Marshal(response)
		_, err = conn.Write(responseBytes)
		socket.CheckError(err)
		fmt.Printf("Write response %s\n", string(responseBytes))
	}
}

// 接收多个客户端请求
func main() {
	// 设置TCP端点的地址
	ip := "127.0.0.1" // ip也可设置成 0.0.0.0 和 空字符串
	port := 5656      // 改成 1023，会报错 bind: permission denied
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(port))
	socket.CheckError(err)

	// 设置监听地址
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	socket.CheckError(err)
	fmt.Println("Waiting for client connection ...")

	// 循环等待TCP连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("Establish connection to client %s\n", conn.RemoteAddr().String()) // 操作系统会随机给客户端分配一个 49152 z~ 65535 上的端口号
		go handleLongRequest(conn)
	}
}
