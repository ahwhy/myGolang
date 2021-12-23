package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ahwhy/myGolang/network/socket"
)

// 收发简单的字符串消息
func TCPServer_SimpleMessage() {
	// 设置TCP端点地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socket.IP+":"+strconv.Itoa(socket.Port))
	socket.CheckError(err)

	// 设置监听地址
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	socket.CheckError(err)
	fmt.Println("Waiting for client connection ...")

	// 等待TCP连接
	conn, err := listener.Accept()
	socket.CheckError(err)
	fmt.Printf("Establish connection to client %s\n", conn.RemoteAddr().String()) // 操作系统会随机给客户端分配一个 49152 ~ 65535 上的端口号

	// 获取客户端Request
	request := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	_, err = conn.Read(request)
	socket.CheckError(err)
	fmt.Printf("Receive request %s\n", string(request)) // []byte转string时，0后面的会自动被截掉

	// 返回Response
	msg, err := conn.Write([]byte("Holle" + string(request)))
	socket.CheckError(err)
	fmt.Printf("Write response %d bytes\n", msg)

	// 中断TCP连接后尝试写入
	// time.Sleep(1 * time.Second)
	// _, err = conn.Write([]byte("Oops"))
	// socket.CheckError(err) // client端已关闭连接，再往conn里写会发生错误

	// 中断TCP连接
	conn.Close()
}
