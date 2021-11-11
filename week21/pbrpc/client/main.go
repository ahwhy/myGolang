package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/ahwhy/myGolang/week21/pbrpc/codec/client"
	"github.com/ahwhy/myGolang/week21/pbrpc/service"
)

// 约束客户端
var _ service.HelloService = (*HelloServiceClient)(nil)

type HelloServiceClient struct {
	*rpc.Client
}

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	// 建立链接
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	// 采用Json编解码的客户端
	c := rpc.NewClientWithCodec(client.NewClientCodec(conn))
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(req *service.Request, resp *service.Response) error {
	// 通过client.Call调用具体的RPC方法
	// 在调用client.Call时
	//   第一个参数是用点号链接的RPC服务名字和方法名字
	//   第二个参数是 请求参数
	//   第三个是请求响应，必须是一个指针，由底层rpc服务进行赋值
	return p.Client.Call(service.HelloServiceName+".Hello", req, resp)
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	resp := &service.Response{}
	err = client.Hello(&service.Request{Value: "hello"}, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
