package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/ahwhy/myGolang/historys/week11/socket"
)

// 设置建立UDP连接数量
const P = 10

// 建立UDP长连接
func main() {
	// 设置UDP端点的地址并向服务端建立UDP连接
	ip := "127.0.0.1" // ip也可设置成 0.0.0.0 和 空字符串
	port := 5656      // 改成 1023，会报错 bind: permission denied
	//跟tcp_client的唯一区别
	conn, err := net.DialTimeout("udp", ip+":"+strconv.Itoa(port), 30*time.Minute) // 30s 后连接超时，一个conn绑定一个本地端口
	socket.CheckError(err)
	fmt.Printf("LocalAddr: %s\nEstablish connection to server %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())

	// 中断UDP连接
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		// 向服务端发送Request

		//多协程，共用一个conn
		go func() {
			defer wg.Done()

			for { //长连接，即连接建立后进行多轮的读写交互
				request := socket.Request{A: 7, B: 5}
				requestBytes, _ := json.Marshal(request)
				_, err = conn.Write(requestBytes)
				socket.CheckError(err)
				fmt.Printf("Write request %s bytes\n", string(requestBytes))

				// 接收服务端发送Response
				responseBytes := make([]byte, 256) // 设定一个最大长度，防止 flood attack；初始化后byte数组每个元素都是0
				read_len, err := conn.Read(responseBytes)
				socket.CheckError(err)
				var response socket.Response
				json.Unmarshal(responseBytes[:read_len], &response) // json反序列化时会把0都考虑在内，需要指定只读前read_len个字节
				fmt.Printf("Receive response: %v\n", response.Sum)
				time.Sleep(1 * time.Second)
			}
		}()
	}
	wg.Wait()
}
