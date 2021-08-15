package main

import (
	"fmt"
	"github.com/ahwhy/myGolang/week11/socket"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader
}

func NewWsServer(port int) *WsServer {
	ws := new(WsServer)
	ws.addr = "0.0.0.0:" + strconv.Itoa(port)
	ws.upgrade = &websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second, //握手超时时间
		ReadBufferSize:   2048,            //读缓冲大小
		WriteBufferSize:  1024,            //写缓冲大小
		//请求检查函数，用于统一的链接检查，以防止跨站点请求伪造。如果Origin请求头存在且原始主机不等于请求主机头，则返回false
		CheckOrigin: func(r *http.Request) bool {
			fmt.Printf("request url %s\n", r.URL)
			fmt.Println("handshake request header")
			for key, values := range r.Header {
				fmt.Printf("%s:%s\n", key, values[0])
			}
			return true
		},
		//http错误响应函数
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {},
	}
	return ws
}

//httpHandler必须实现ServeHTTP接口
func (ws *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		httpCode := http.StatusInternalServerError //返回一个内部错误的信息
		rePhrase := http.StatusText(httpCode)      //组织错误话术
		fmt.Println("path error ", rePhrase)
		http.Error(w, rePhrase, httpCode) //把出错的话术写到ResponseWriter里
		return
	}
	conn, err := ws.upgrade.Upgrade(w, r, nil) //将http协议升级到websocket协议
	if err != nil {
		fmt.Printf("upgrade http to websocket error: %v\n", err)
		return
	}
	fmt.Printf("establish conection to client %s\n", conn.RemoteAddr().String())
	go ws.handleConnection(conn)
}

//处理连接里发来的请求数据
func (ws *WsServer) handleConnection(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()
	for { //长连接
		conn.SetReadDeadline(time.Now().Add(20 * time.Second))
		var request socket.Request
		if err := conn.ReadJSON(&request); err != nil {
			//判断是不是超时
			if netError, ok := err.(net.Error); ok { //如果ok==true，说明类型断言成功
				if netError.Timeout() {
					fmt.Printf("read message timeout, remote %s\n", conn.RemoteAddr().String())
					return
				}
			}
			//忽略websocket.CloseGoingAway/websocket.CloseNormalClosure这2种closeErr，如果是其他closeErr就打一条错误日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Printf("read message from %s error %v\n", conn.RemoteAddr().String(), err)
			}
			return //只要ReadMessage发生错误，就关闭这条连接
		} else {
			response := socket.Response{Sum: request.A + request.B}
			if err = conn.WriteJSON(&response); err != nil {
				fmt.Printf("write response failed: %v", err)
			} else {
				fmt.Printf("write response %d\n", response.Sum)
			}
		}
	}
}

func (ws *WsServer) Start() (err error) {
	ws.listener, err = net.Listen("tcp", ws.addr) //http和websocket都是建立在tcp之上的
	if err != nil {
		fmt.Printf("listen error:%s\n", err)
		return
	}
	err = http.Serve(ws.listener, ws) //开始对外提供http服务。可以接收很多连接请求，其他一个连接处理出错了，也不会影响其他连接
	if err != nil {
		fmt.Printf("http server error: %v\n", err)
		return
	}

	// if err:=http.ListenAndServe(ws.addr, ws);err!=nil{	//Listen和Serve两步合成一步
	// 	fmt.Printf("http server error: %v\n", err)
	// 	return
	// }
	return nil
}

func main_ws_server() {
	ws := NewWsServer(5657)
	ws.Start()
}

//go run socket/server/ws_server.go
