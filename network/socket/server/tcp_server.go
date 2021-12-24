package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

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
	requestBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	_, err = conn.Read(requestBytes)
	socket.CheckError(err)
	fmt.Printf("Receive request %s\n", string(requestBytes)) // []byte转string时，0后面的会自动被截掉

	// 返回Response
	write_len, err := conn.Write([]byte("Holle " + string(requestBytes)))
	socket.CheckError(err)
	fmt.Printf("Write response %d bytes\n", write_len)

	// 中断TCP连接后尝试写入
	// time.Sleep(1 * time.Second)
	// _, err = conn.Write([]byte("Oops"))
	// socket.CheckError(err) // client端已关闭连接，再往conn里写会发生错误

	// 中断TCP连接
	conn.Close()
}

// 序列化请求和响应的结构体
func TCPServer_StructMessage() {
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
	requestBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	read_len, err := conn.Read(requestBytes)
	socket.CheckError(err)
	fmt.Printf("Receive request %s\n", string(requestBytes)) // []byte转string时，0后面的会自动被截掉

	// 返回Response
	var request socket.Request
	json.Unmarshal(requestBytes[:read_len], &request) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节)
	response := socket.Response{Sum: request.A + request.B}
	responseBytes, _ := json.Marshal(response)
	_, err = conn.Write(responseBytes)
	socket.CheckError(err)
	fmt.Printf("Write response %s\n", string(responseBytes))

	// 中断TCP连接
	conn.Close()
}

// 处理多个客户端请求
func TCPServer_MoreStructMessage() {
	// 设置TCP端点地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socket.IP+":"+strconv.Itoa(socket.Port))
	socket.CheckError(err)

	// 设置监听地址
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	socket.CheckError(err)
	fmt.Println("Waiting for client connection ...")

	// 循环等待TCP连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			socket.CheckError(err)
			continue
		}
		fmt.Printf("Establish connection to client %s\n", conn.RemoteAddr().String()) // 操作系统会随机给客户端分配一个 49152 ~ 65535 上的端口号
		go handleRequest(conn)
	}
}

// 处理Request
func handleRequest(conn net.Conn) {
	// 中断TCP连接
	defer conn.Close()

	// 获取客户端Request
	requestBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	read_len, err := conn.Read(requestBytes)
	socket.CheckError(err)
	fmt.Printf("Receive request %s\n", string(requestBytes)) // []byte转string时，0后面的会自动被截掉

	// 返回Response
	var request socket.Request
	json.Unmarshal(requestBytes[:read_len], &request) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节)
	response := socket.Response{Sum: request.A + request.B}
	responseBytes, _ := json.Marshal(response)
	_, err = conn.Write(responseBytes)
	socket.CheckError(err)
	fmt.Printf("Write response %s\n", string(responseBytes))
}

// 处理多个客户端请求
func TCPServer_MoreLongStructMessage() {
	// 设置TCP端点地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socket.IP+":"+strconv.Itoa(socket.Port))
	socket.CheckError(err)

	// 设置监听地址
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	socket.CheckError(err)
	fmt.Println("Waiting for client connection ...")

	// 循环等待TCP连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			socket.CheckError(err)
			continue
		}
		fmt.Printf("Establish connection to client %s\n", conn.RemoteAddr().String()) // 操作系统会随机给客户端分配一个 49152 ~ 65535 上的端口号
		go handleLongRequest(conn)
	}
}

// 处理Request
func handleLongRequest(conn net.Conn) {
	// 设置超时时间，30秒后 conn.Read 会报出 i/o timeout
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	// 中断TCP连接
	defer conn.Close()

	// 长连接，即连接建立后进行多轮的读写互交
	for {
		// 获取客户端Request
		requestBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
		read_len, err := conn.Read(requestBytes)
		if err != nil {
			socket.CheckError(err)
			break // 到达deadline后，退出for循环，关闭连接；client再用这个连接读写会发生错误
		}
		fmt.Printf("Receive request %s\n", string(requestBytes)) // []byte转string时，0后面的会自动被截掉

		// 返回Response
		var request socket.Request
		json.Unmarshal(requestBytes[:read_len], &request) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节)
		response := socket.Response{Sum: request.A + request.B}
		responseBytes, _ := json.Marshal(response)
		_, err = conn.Write(responseBytes)
		socket.CheckError(err)
		fmt.Printf("Write response %s\n", string(responseBytes))
	}
}
