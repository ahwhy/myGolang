package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second //
	pongWait   = 60 * time.Second //
	pingPeriod = 9 * pongWait / 10
	maxMsgSize = 512 //消息的长度不能超过512
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 2 * time.Second, //握手超时时间
	ReadBufferSize:   1024,            //读缓冲大小
	WriteBufferSize:  1024,            //写缓冲大小
	CheckOrigin:      func(r *http.Request) bool { return true },
	Error:            func(w http.ResponseWriter, r *http.Request, status int, reason error) {},
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	frontName []byte //前端的名字，用于展示在消息前面
}

//从websocket连接里读出数据，发给hub
func (client *Client) read() {
	defer func() { //收尾工作
		client.hub.unregister <- client //从hub那注销client
		fmt.Printf("close connection to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close() //关闭websocket管道
	}()
	client.conn.SetReadLimit(maxMsgSize) //一次从管管中读取的最大长度
	/**
	连接不断的情况下，每隔54秒向客户端发一次ping，客户端返回pong，所以把ReadDeadline设为60秒是没有问题的
	*/
	client.conn.SetReadDeadline(time.Now().Add(pongWait)) //60秒后不允许读
	client.conn.SetPongHandler(func(appData string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait)) //每次收到pong都把deadline往后推迟60秒
		return nil
	})
	for {
		_, message, err := client.conn.ReadMessage() //如果前端主动断开连接，该行会报错，for循环会退出。注册client时，hub那儿会关闭client.send管道
		if err != nil {
			//如果以意料之外的关闭状态关闭，就打印日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				fmt.Printf("close websocket conn error: %v\n", err)
			}
			break //只要ReadMessage失败，就关闭websocket管道、那注销client，退出
		} else {
			//换行符用空格替代，bytes.TrimSpace把首尾连续的空格去掉
			message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
			if len(client.frontName) == 0 {
				client.frontName = message //赋给frontName，不进行广播
				fmt.Printf("%s online\n", string(client.frontName))
			} else {
				//要广播的内容前面加上front的名字
				client.hub.broadcast <- bytes.Join([][]byte{client.frontName, message}, []byte(": ")) //从websocket连接里读出数据，发给hub的broadcast
			}
		}
	}
}

//从hub的broadcast那儿读限数据，写到websocket连接里面去
func (client *Client) write() {
	ticker := time.NewTicker(pingPeriod) //给前端发心跳，看前端是否还存活
	defer func() {
		ticker.Stop() //ticker不用就stop，防止协程泄漏
		fmt.Printf("close connection to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close() //给前端写数据失败，就可以关系连接了
	}()

	for {
		select {
		case msg, ok := <-client.send: //正常情况是hub发来了数据。如果前端断开了连接，read()会触发client.send管道的关闭，该case会立即执行。从而执行!ok里的return，从而执行defer
			if !ok { //client.send该管道被hub关闭了
				client.conn.WriteMessage(websocket.CloseMessage, []byte{}) //写一条关闭信息就可以结束一切了
				return
			}
			//思考：如果把SetWriteDeadline这行代码放到for循环上面会怎样？向conn里写数据就有可能报i/o timeout
			client.conn.SetWriteDeadline(time.Now().Add(writeWait)) //10秒内必须把信息写给前端（写到websocket连接里去），否则就关闭连接
			/**
			消息类型有5种：TextMessage，BinaryMessage，CloseMessage，PingMessage，PongMessage
			*/
			if writer, err := client.conn.NextWriter(websocket.TextMessage); err != nil { //通过NextWriter创建一个新的writer，主要是为了确保上一个writer已经被关闭，即它想写的内容已经flush到conn里去了
				return
			} else {
				writer.Write(msg)
				writer.Write(newLine) //每发一条消息，都加一个换行符
				//为了提升性能，如果client.send里还有消息，则趁这一次都写给前端
				n := len(client.send)
				for i := 0; i < n; i++ {
					writer.Write(<-client.send)
					writer.Write(newLine)
				}
				if err := writer.Close(); err != nil {
					return //结束一切
				}
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			//心跳保持，给浏览器发一个PingMessage，等待浏览器返回PongMessage
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return //写websocket连接失败，说明连接出问题了，该client可以over了
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) //http升级为websocket协议
	if err != nil {
		fmt.Printf("upgrade error: %v\n", err)
		return
	}
	fmt.Printf("connect to client %s\n", conn.RemoteAddr().String())
	//每来一个前端请求，就会创建一个client
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//向hub注册client
	client.hub.register <- client

	//启动子协程，运行ServeWs的协程退出后子协程也不会能出
	//websocket是全双工模式，可以同时read和write
	go client.read()
	go client.write()
}
