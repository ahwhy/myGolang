package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/network/socket"
)

// 处理多个客户端请求
func UDPServer_MoreLongStructMessage() {
	// 设置UDP端点地址
	udpAddr, err := net.ResolveUDPAddr("udp4", socket.IP+":"+strconv.Itoa(socket.Port))
	socket.CheckError(err)

	// 设置UDP监听地址并等待连接
	conn, err := net.ListenUDP("udp4", udpAddr) // UDP不需要创建连接，不需要像TCP那样通过 Accept()创建连接，这里的conn是个假连接
	socket.CheckError(err)

	// 设置30s超时时间
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	// 中断UDP连接
	defer conn.Close()

	//长连接，即连接建立后进行多轮的读写交互
	for {
		// 获取客户端Request
		requestBytes := make([]byte, 256)                           // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
		read_len, remoteAddr, err := conn.ReadFromUDP(requestBytes) // 一个UDPconn可以对应多个client
		if err != nil {
			socket.CheckError(err)
			break // 到达deadline后，退出for循环，关闭连接；client再用这个连接读写会发生错误
		}
		fmt.Printf("Receive request %s from %s\n", string(requestBytes), remoteAddr.String()) // []byte转string时，0后面的会自动被截掉

		// 返回Response
		var request socket.Request
		json.Unmarshal(requestBytes[:read_len], &request) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节)
		response := socket.Response{Sum: request.A + request.B}
		responseBytes, _ := json.Marshal(response)
		_, err = conn.WriteToUDP(responseBytes, remoteAddr)
		socket.CheckError(err)
		fmt.Printf("Write response %s to %s\n", string(responseBytes), remoteAddr.String())
	}
}
