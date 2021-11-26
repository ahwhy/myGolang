package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string // 对外发送消息的通道

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // 所有接受的客户消息
)

func broadcaster() {
	clients := make(map[client]bool) // 所有连接的客户端集合
	for {
		select {
		case msg := <-messages:
			// 把所有接收的消息广播给所有的客户
			// 发送消息通道
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // 对外发送客户消息的通道
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who           // 这条单发给自己
	messages <- who + " has arrived" // 这条进行进行广播，但是自己还没加到广播列表中
	entering <- ch                   // 然后把自己加到广播列表中

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// 注意，忽略input.Err()中可能的错误

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// 客户端处理函数: clientWriter
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		// 在消息结尾使用 \r\n ，提升平台兼容
		fmt.Fprintf(conn, "%s\r\n", msg) // 注意，忽略网络层面的错误
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
