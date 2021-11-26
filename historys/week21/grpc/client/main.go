package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ahwhy/myGolang/week21/grpc/service"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := service.NewHelloServiceClient(conn)

    // 通过接口定义的方法就可以调用服务端对应的gRPC服务提供的方法
	// req := &service.Request{Value: "hello"}
	// reply, err := client.Hello(context.Background(), req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(reply.GetValue())

	// 客户端需要先调用Channel方法获取返回的流对象
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 客户端将发送和接收操作放到两个独立的Goroutine
	// 首先是向服务端发送数据
	go func() {
		for {
			if err := stream.Send(&service.Request{Value: "hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	// 然后在循环中接收服务端返回的数据
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
}
