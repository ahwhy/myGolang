package protocol

import (
	"log"
	"net/http"
	"time"

	"github.com/ahwhy/myGolang/network/chatroom/module"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 2 * time.Second, // 握手超时时间
	ReadBufferSize:   1024,            // 读缓冲大小
	WriteBufferSize:  1024,            // 写缓冲大小
	CheckOrigin:      func(r *http.Request) bool { return true },
	Error:            func(w http.ResponseWriter, r *http.Request, status int, reason error) {},
}

func ServeWs(hub *module.Hub, w http.ResponseWriter, r *http.Request) {
	// 升级http为websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v\n", err)
		return
	}
	log.Printf("connect to client %s\n", conn.RemoteAddr().String())
	// 每来一个前端请求，就会创建一个client
	client := module.NewClient(conn, hub)
	client.Hub.Register <- client

	// 启动子协程，运行ServeWs的协程退出后子协程也不会能出
	// websocket是全双工模式，可以同时read和write
	go client.Read()
	go client.Write()
}
