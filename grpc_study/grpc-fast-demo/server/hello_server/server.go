package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc-fast-demo/proto/hello"
	"log"
	"net"
	"time"
)

type HelloService struct {
	// 必须嵌入UnimplementedUserServiceServer
	hello.UnimplementedUserServiceServer
}

// 实现Say方法
func (h *HelloService) Say(ctx context.Context, req *hello.Request) (res *hello.Response, err error) {
	format := time.Now().Format("2006-01-02 15:04:05")
	return &hello.Response{Result: "你好 " + req.GetName() + "---" + format}, nil
}

const PORT = "8888"

func main() {
	// 创建grpc服务
	server := grpc.NewServer()
	// 注册服务
	hello.RegisterUserServiceServer(server, &HelloService{})

	// 监听端口
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}
