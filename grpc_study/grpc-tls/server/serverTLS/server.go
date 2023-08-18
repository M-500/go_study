package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-tls/pb/tls_demo"
	"log"
	"net"
	"time"
)

const port = ":8999"

type TLSService struct {
	tls_demo.UnimplementedHelloServiceServer
}

func (t *TLSService) SayHello(ctx context.Context, in *tls_demo.HelloRequest) (*tls_demo.HelloResponse, error) {
	format := time.Now().Format("2006-01-02 15:04:05")
	return &tls_demo.HelloResponse{Res: "您好 " + in.GetName() + "---" + format}, nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}

	cred, err := credentials.NewServerTLSFromFile("./conf/server.pem", "./conf/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(cred))
	tls_demo.RegisterHelloServiceServer(s, &TLSService{})

	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("close err:%v", err)
		return
	}

}
