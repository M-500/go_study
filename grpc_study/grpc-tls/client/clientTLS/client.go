package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-tls/pb/tls_demo"
	"log"
)

const port = ":8999"

func main() {
	// 根据客户端输入的证书文件和密钥构造 TLS 凭证。
	// 第二个参数 serverNameOverride 为服务名称。
	c, err := credentials.NewClientTLSFromFile("./conf/server.pem", "grpc-tls")

	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}
	// 返回一个配置连接的 DialOption 选项。
	// 用于 grpc.Dial(target string, opts ...DialOption) 设置连接选项
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := tls_demo.NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &tls_demo.HelloRequest{
		Name: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetRes())
}
