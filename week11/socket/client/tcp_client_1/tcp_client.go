package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ahwhy/myGolang/week11/socket"
)

// 收发简单的字符串消息
func main() {
	// 设置TCP端点的地址
	ip := "127.0.0.1"                                                     // ip也可设置成 0.0.0.0 和 空字符串
	port := 5656                                                          // 改成 1023，会报错 bind: permission denied
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(port)) // 也可设置 www.baidu.com:80
	socket.CheckError(err)
	fmt.Printf("Ip: %s Port: %d\n", tcpAddr.IP.String(), tcpAddr.Port)

	// 向服务端建立TCP连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	socket.CheckError(err)
	fmt.Printf("Establish connection to server %s\n", conn.RemoteAddr().String())

	// 向服务端发送Request
	n, err := conn.Write([]byte("World"))
	socket.CheckError(err)
	fmt.Printf("Write request %d bytes\n", n)

	// 接收服务端发送Response
	response := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	_, err = conn.Read(response)
	// response, err := ioutil.ReadAll(conn) //从conn中读取所有内容，直到遇到error(比如连接关闭)或EOF
	socket.CheckError(err)
	fmt.Printf("Receive response: %s\n", string(response)) // []byte转string时，0后面的会自动被截掉

	// 中断TCP连接
	conn.Close()
}
