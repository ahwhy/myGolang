package client

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ahwhy/myGolang/network/socket"
)

// 收发简单的字符串消息
func TCPClient_SimpleMessage() {
	// 设置TCP端点地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socket.IP+":"+strconv.Itoa(socket.Port)) // 也可设置 www.baidu.com:80
	socket.CheckError(err)
	fmt.Printf("IP: %s Port: %d\n", tcpAddr.IP.String(), tcpAddr.Port)

	// 向服务端建立TCP连接
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	socket.CheckError(err)
	fmt.Printf("Establish connection to server %s\n", conn.RemoteAddr().String())

	// 向服务端发送Request
	msg, err := conn.Write([]byte("World"))
	socket.CheckError(err)
	fmt.Printf("Write request %d bytes\n", msg)

	// 接收服务端发送的Response
	response := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	_, err = conn.Read(response)
	// response, err := ioutil.ReadAll(conn) // 从conn中读取所有内容，直到遇到error(比如连接关闭)或EOF
	socket.CheckError(err)
	fmt.Printf("Receive response: %s\n", string(response)) // []byte转string时，0后面的会自动被截掉

	// 中断TCP连接
	conn.Close()
}
