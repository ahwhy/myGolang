package module

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	FrontName []byte // 前端的名字，用于展示在消息前面
}

// 从websocket连接里读出数据，发给hub
func (client *Client) Read() {
	defer func() {
		client.Hub.Unregister <- client // 从hub那注销client
		log.Printf("close connection to %s\n", client.Conn.RemoteAddr().String())
		client.Conn.Close() //关闭websocket管道
	}()

	client.Conn.SetReadLimit(maxMsgSize)                  // 一次从管道中读取的最大长度
	client.Conn.SetReadDeadline(time.Now().Add(pongWait)) // 60秒后不允许读
	// 连接不断的情况下，每隔54秒向客户端发一次ping，客户端返回pong，所以把ReadDeadline设为60秒是没有问题的
	client.Conn.SetPongHandler(func(appData string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait)) // 每次收到pong都把deadline往后推迟60秒
		return nil
	})

	for {
		// 如果前端主动断开连接，该行会报错，for循环会退出；注销client后，hub会关闭client.send管道
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("close websocket conn error: %v\n", err)
			}
			break // 只要ReadMessage失败，就关闭websocket管道、注销client，退出
		} else {
			// 换行符用空格替代，bytes.TrimSpace把首尾连续的空格去掉
			message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
			if len(client.FrontName) == 0 {
				client.FrontName = message // 第一条消息赋给frontName，并且不进行广播
				log.Printf("%s online\n", string(client.FrontName))
			} else {
				// 从websocket连接里读出数据，发给hub的broadcast并且广播的内容前面加上frontName，//
				client.Hub.Broadcast <- bytes.Join([][]byte{client.FrontName, message}, []byte(": "))
			}
		}
	}
}

// 从hub的broadcast 读出数据，写到websocket连接里面去
func (client *Client) Write() {
	// 设置给前端发的心跳
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		log.Printf("close connection to %s\n", client.Conn.RemoteAddr().String())
		client.Conn.Close()
	}()

	for {
		select {
		// 正常情况下 hub发来数据后，如果前端断开连接，read()会触发client.send管道的关闭，该case会立即执行，即 执行!ok里的return，最后执行defer
		case msg, ok := <-client.Send:
			if !ok { // client.send 管道被hub关闭
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{}) // 发送一条关闭信息后，就结束一切
				return
			}
			// 如果把SetWriteDeadline这行代码放到for循环上面，向conn里写数据就有可能报i/o timeout
			// 10秒内必须把信息写给前端(写到websocket连接里去)，否则就关闭连接
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 消息类型有5种：TextMessage，BinaryMessage，CloseMessage，PingMessage，PongMessage
			// 通过NextWriter创建一个新的writer，主要是为了确保上一个writer已经被关闭，即它想写的内容已经flush到conn里
			if writer, err := client.Conn.NextWriter(websocket.TextMessage); err != nil {
				return
			} else {
				writer.Write(msg)
				writer.Write(newLine)
				// 为了提升性能，如果client.send里还有消息，则趁这一次都写给前端
				n := len(client.Send)
				for i := 0; i < n; i++ {
					writer.Write(<-client.Send)
					writer.Write(newLine)
				}
				if err := writer.Close(); err != nil {
					return
				}
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 心跳保持，给浏览器发一个PingMessage，等待浏览器返回PongMessage
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return // 写websocket连接失败，说明连接出现问题，该client可以over
			}
		}
	}
}
