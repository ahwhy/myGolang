package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/ahwhy/myGolang/week21/pb"
)

func main() {
	clientObj := &pb.String{Value: "hello proto3"}

	// of := &pb.SampleMessage{}
	// of.GetSub1()
	// of.GetSub2()

	// 序列化
	out, err := proto.Marshal(clientObj)
	if err != nil {
		log.Fatalln("Failed to encode obj:", err)
	}

	// 二进制编码
	fmt.Println("encode bytes: ", out)

	// 反序列化
	serverObj := &pb.String{}
	err = proto.Unmarshal(out, serverObj)
	if err != nil {
		log.Fatalln("Failed to decode obj:", err)
	}
	fmt.Println("decode obj: ", serverObj)
}
