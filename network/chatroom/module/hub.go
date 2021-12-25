package module

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte), //同步管道，确保hub收到的消息不会堆积；如果同时有多个client想给hub发数据就阻塞
	}
}

type Hub struct {
	Clients    map[*Client]bool // 维护所有Client
	Register   chan *Client     // Client注册请求通过管道来接收
	Unregister chan *Client     // Client注销请求通过管道来接收
	Broadcast  chan []byte      // 需要广播的消息
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true // 注册client
		case client := <-hub.Unregister:
			if _, ok := hub.Clients[client]; ok { // 防止重复注销
				delete(hub.Clients, client) // 注销client
				close(client.Send)          // hub从此以后不需要再向该client广播消息
			}
		case msg := <-hub.Broadcast:
			for client := range hub.Clients {
				select {
				case client.Send <- msg: // 如果管道不能立即写入数据，就认为该client出现故障
				default:
					close(client.Send)
					delete(hub.Clients, client)
				}
			}
		}
	}
}
