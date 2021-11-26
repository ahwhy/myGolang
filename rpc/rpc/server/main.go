package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/ahwhy/myGolang/rpc/rpc/service"
)

// 通过接口约束HelloService服务
var _ service.HelloService = (*HelloService)(nil)

type HelloService struct{}

// Hello的逻辑 将对方发送的消息前面添加一个Hello 然后返还给对方
// 由于是一个rpc服务，参数上的约束为
//   第一个参数是请求
//   第二个参数是响应
// 可以类比Http handler
func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
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
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// - 新的RPC服务其实是一个类似REST规范的接口，接收请求并采用相应处理流程
// - 首先依然要解决JSON编解码的问题，需要将HTTP接口的Handler参数传递给jsonrpc，
// 	- 满足jsonrpc接口，提前构建io.ReadWriteCloser类型的conn通道
// 	- Writer是ResponseWriter，ReadCloser是Request的Body，直接内嵌就可以
// func NewRPCReadWriteCloserFromHTTP(w http.ResponseWriter, r *http.Request) *RPCReadWriteCloser {
// 	return &RPCReadWriteCloser{w, r.Body}
// }

// type RPCReadWriteCloser struct {
// 	io.Writer
// 	io.ReadCloser
// }

// func main() {
// 	rpc.RegisterName("HelloService", new(HelloService))

// RPC的服务架设在"/jsonrpc"路径
// 在处理函数中基于http.ResponseWriter和http.Request类型的参数构造一个io.ReadWriteCloser类型的conn通道
// 然后基于conn构建针对服务端的json编码解码器
// 最后通过rpc.ServeRequest函数为每次请求处理一次RPC方法调用
// 	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
// 		conn := NewRPCReadWriteCloserFromHTTP(w, r)
// 		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
// 	})

// 	http.ListenAndServe(":1234", nil)
// }
