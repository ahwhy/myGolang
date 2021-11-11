package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/ahwhy/myGolang/week21/pbrpc/codec/server"
	"github.com/ahwhy/myGolang/week21/pbrpc/service"
)

// 通过接口约束HelloService服务
var _ service.HelloService = (*HelloService)(nil)

type HelloService struct{}

// Hello的逻辑 将对方发送的消息前面添加一个Hello 然后返还给对方
// 由于是一个rpc服务，参数上的约束为
//   第一个参数是请求
//   第二个参数是响应
// 可以类比Http handler
func (p *HelloService) Hello(req *service.Request, resp *service.Response) error {
	resp.Value = "hello:" + req.Value
	return nil
}

func main() {
	// 把HelloService对象注册成一个rpc的 receiver
	// 其中rpc.Register函数调用会将对象类型中所有满足RPC规则的对象方法注册为RPC函数，
	// 所有注册的方法会放在"HelloService"服务空间之下
	rpc.RegisterName("HelloService", new(HelloService))

	// 然后建立一个唯一的TCP链接
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	// 通过rpc.ServeConn函数在该TCP链接上为对方提供RPC服务
	// 每Accept一个请求，就创建一个goroutie进行处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// 前面都是tcp的知识，这里之后就是RPC接管
		// 因此可以认为 rpc 封装消息到函数调用的这个逻辑
		// 提升了工作效率，逻辑比较简洁
		// go rpc.ServeConn(conn)

		// 代码中最大的变化是用rpc.ServeCodec函数替代了rpc.ServeConn函数
		// 传入的参数是针对服务端的json编解
		go rpc.ServeCodec(server.NewServerCodec(conn))
	}
}
