package main

import (
	"context"
	"io"
	"log"
	"net"

	"gitee.com/infraboard/go-course/day21/grpc/service"
	"google.golang.org/grpc"
)

var _ service.HelloServiceServer = (*HelloService)(nil)

type HelloService struct {
	service.UnimplementedHelloServiceServer
}

func (p *HelloService) Hello(ctx context.Context, req *service.Request) (*service.Response, error) {
	resp := &service.Response{}
	resp.Value = "hello:" + req.Value
	return resp, nil
}

func (p *HelloService) Channel(stream service.HelloService_ChannelServer) error {
	for {
		// 接收一个请求
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// 响应一个请求
		resp := &service.Response{Value: "hello:" + args.GetValue()}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
}

func main() {
	// 首先是通过grpc.NewServer()构造一个gRPC服务对象
	grpcServer := grpc.NewServer()
	// 然后通过gRPC插件生成的RegisterHelloServiceServer函数注册我们实现的HelloServiceImpl服务
	service.RegisterHelloServiceServer(grpcServer, new(HelloService))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	// 然后通过grpcServer.Serve(lis)在一个监听端口上提供gRPC服务
	grpcServer.Serve(lis)
}
