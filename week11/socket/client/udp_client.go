package main

import (
	"encoding/json"
	"fmt"
	"github.com/ahwhy/myGolang/week11/socket"
	"net"
	"strconv"
	"sync"
	"time"
)

//长连接
func main_udp_client() {
	ip := "127.0.0.1" //ip换成0.0.0.0和空字符串试试
	port := 5656
	//跟tcp_client的唯一区别就是这行代码
	conn, err := net.DialTimeout("udp", ip+":"+strconv.Itoa(port), 30*time.Minute) //一个conn绑定一个本地端口
	socket.CheckError(err)
	defer conn.Close()
	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		request := socket.Request{A: 7, B: 4}
		requestBytes, _ := json.Marshal(request)
		go func() { //多协程，共用一个conn
			defer wg.Done()
			for { //长连接，即连接建立后进行多轮的读写交互
				_, err = conn.Write(requestBytes)
				socket.CheckError(err)
				fmt.Printf("write request %s\n", string(requestBytes))
				responseBytes := make([]byte, 256) //初始化后byte数组每个元素都是0
				read_len, err := conn.Read(responseBytes)
				socket.CheckError(err)
				var response socket.Response
				json.Unmarshal(responseBytes[:read_len], &response) //json反序列化时会把0都考虑在内，所以需要指定只读前read_len个字节
				fmt.Printf("receive response: %d\n", response.Sum)
				time.Sleep(1 * time.Second)
			}
		}()
	}
	wg.Wait()
}

//先启动tcp_server，再在另外一个终端启动udp_client
//go run socket/client/udp_client.go
