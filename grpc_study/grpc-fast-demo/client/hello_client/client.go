package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-fast-demo/proto/hello"
	"log"
)

const PORT = "8888"

func main() {
	// 建立链接
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	// 一定要记得关闭链接
	defer conn.Close()

	// 实例化客户端
	client := hello.NewUserServiceClient(conn)
	// 发起请求
	response, err := client.Say(context.Background(), &hello.Request{Name: "gRPC测试客户端"})
	if err != nil {
		log.Fatalf("client.Say err: %v", err)
	}
	fmt.Printf("返回值为: %s", response.GetResult())

}
