package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/network/socket"
	"github.com/gorilla/websocket"
)

func NewWsClient(ip string, port int) *WsClient {
	return &WsClient{
		addr: "ws://" + ip + ":" + strconv.Itoa(port),
		// 构造Request报文头部 Cookie:name=atlantis
		header: &http.Header{
			"Cookie": []string{"name=Atlantis"},
		},
		dialer: &websocket.Dialer{},
	}
}

type WsClient struct {
	addr   string
	header *http.Header
	dialer *websocket.Dialer
}

func (ws *WsClient) Start() {
	// 建立websocket连接
	conn, resp, err := ws.dialer.Dial(ws.addr+"/add", *ws.header)
	if err != nil {
		fmt.Printf("Dial server error:%v\n", err)
		return
	}
	defer conn.Close()

	// 打印Response报文头部
	fmt.Println("Handshake response header")
	for key, values := range resp.Header {
		fmt.Printf("%s:%s\n", key, values[0])
	}

	for i := 0; i < P; i++ {
		// 发送Request
		request := socket.Request{A: 7, B: 4}
		err = conn.WriteJSON(request) //websocket.Conn直接提供发json序列化和反序列化方法
		socket.CheckError(err)
		requestBytes, _ := json.Marshal(request)
		fmt.Printf("write request %s\n", string(requestBytes))

		// 接收服务端发送Response
		var response socket.Response
		err = conn.ReadJSON(&response)
		socket.CheckError(err)
		fmt.Printf("receive response: %d\n", response.Sum)
		time.Sleep(1 * time.Second)
	}
}

func WSClient_MoreLongStructMessage() {
	ws := NewWsClient(socket.IP, socket.Port)
	ws.Start()
}
