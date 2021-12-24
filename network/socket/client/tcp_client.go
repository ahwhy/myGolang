package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

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
	write_len, err := conn.Write([]byte("World"))
	socket.CheckError(err)
	fmt.Printf("Write request %d bytes\n", write_len)

	// 接收服务端发送的Response
	responseBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	_, err = conn.Read(responseBytes)
	// response, err := ioutil.ReadAll(conn) // 从conn中读取所有内容，直到遇到error(比如连接关闭)或EOF
	socket.CheckError(err)
	fmt.Printf("Receive response: %s\n", string(responseBytes)) // []byte转string时，0后面的会自动被截掉

	// 中断TCP连接
	conn.Close()
}

// 序列化请求和响应的结构体
func TCPClient_StructMessage() {
	// 设置TCP端点的地址并向服务端建立TCP连接
	conn, err := net.DialTimeout("tcp4", socket.IP+":"+strconv.Itoa(socket.Port), 30*time.Second) // 30s 后连接超时
	socket.CheckError(err)
	fmt.Printf("LocalAddr: %s\nEstablish connection to server %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())

	// 向服务端发送Request
	request := socket.Request{A: 7, B: 5}
	requestBytes, _ := json.Marshal(request)
	_, err = conn.Write(requestBytes)
	socket.CheckError(err)
	fmt.Printf("Write request %s\n", string(requestBytes))

	// 接收服务端发送的Response
	responseBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	read_len, err := conn.Read(responseBytes)
	socket.CheckError(err)
	var response socket.Response
	json.Unmarshal(responseBytes[:read_len], &response) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节
	fmt.Printf("Receive response: %v\n", response.Sum)

	// 中断TCP连接
	conn.Close()
}

// 序列化请求和响应的结构体
func TCPClient_LongStructMessage() {
	// 设置TCP端点的地址并向服务端建立TCP连接
	conn, err := net.DialTimeout("tcp4", socket.IP+":"+strconv.Itoa(socket.Port), 30*time.Second) // 30s 后连接超时
	socket.CheckError(err)
	fmt.Printf("LocalAddr: %s\nEstablish connection to server %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())
	// 中断TCP连接
	defer conn.Close()

	// 长连接，即连接建立后进行多轮的读写交互
	for {
		// 向服务端发送Request
		request := socket.Request{A: 7, B: 5}
		requestBytes, _ := json.Marshal(request)
		_, err = conn.Write(requestBytes)
		socket.CheckError(err)
		fmt.Printf("Write request %s\n", string(requestBytes))

		// 接收服务端发送的Response
		responseBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
		read_len, err := conn.Read(responseBytes)
		socket.CheckError(err)
		var response socket.Response
		json.Unmarshal(responseBytes[:read_len], &response) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节
		fmt.Printf("Receive response: %v\n", response.Sum)
		time.Sleep(1 * time.Second)
	}
}


