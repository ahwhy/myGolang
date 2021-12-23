package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/historys/week11/socket"
)

// 收发简单的字符串消息
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

	// 等待TCP连接
	conn, err := listener.Accept()
	socket.CheckError(err)
	fmt.Printf("Establish connection to client %s\n", conn.RemoteAddr().String()) // 操作系统会随机给客户端分配一个 49152 z~ 65535 上的端口号

	// 获取客户端Request
	request := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	_, err = conn.Read(request)
	socket.CheckError(err)
	fmt.Printf("Receive request %s\n", string(request)) // []byte转string时，0后面的会自动被截掉

	// 返回Response
	n, err := conn.Write([]byte("Hello " + string(request)))
	socket.CheckError(err)
	fmt.Printf("Write response %d bytes\n", n)

	// 中断TCP连接后尝试写入
	time.Sleep(1 * time.Second)
	_, err = conn.Write([]byte("Oops"))
	socket.CheckError(err) //client端已关闭连接，再往conn里写会发生错误
}
