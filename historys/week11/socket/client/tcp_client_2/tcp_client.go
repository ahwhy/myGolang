package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/historys/week11/socket"
)

// 序列化请求和响应的结构体
func main() {
	// 设置TCP端点的地址并向服务端建立TCP连接
	ip := "127.0.0.1"                                                               // ip也可设置成 0.0.0.0 和 空字符串
	port := 5656                                                                    // 改成 1023，会报错 bind: permission denied
	conn, err := net.DialTimeout("tcp4", ip+":"+strconv.Itoa(port), 30*time.Minute) // 30s 后连接超时
	socket.CheckError(err)
	fmt.Printf("LocalAddr: %s\nEstablish connection to server %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())

	// 向服务端发送Request
	request := socket.Request{A: 7, B: 5}
	requestBytes, _ := json.Marshal(request)
	_, err = conn.Write(requestBytes)
	socket.CheckError(err)
	fmt.Printf("Write request %s\n", string(requestBytes))

	// 接收服务端发送Response
	responseBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
	read_len, err := conn.Read(responseBytes)
	socket.CheckError(err)
	var response socket.Response
	json.Unmarshal(responseBytes[:read_len], &response) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节
	fmt.Printf("Receive response: %v\n", response.Sum)

	// 中断TCP连接
	conn.Close()
}

