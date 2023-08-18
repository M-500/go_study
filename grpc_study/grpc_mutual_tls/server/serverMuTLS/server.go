package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_mutual_tls/pb/tls_demo"
	"io/ioutil"
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
	// 公钥中读取和解析公钥/私钥对
	cred, err := tls.LoadX509KeyPair("./conf/server.crt", "./conf/server.key")
	if err != nil {
		log.Fatalf("LoadX509KeyPair  err: %v", err)
	}
	// 创建一组根证书
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		fmt.Println("read ca pem error ", err)
		return
	}

	// 解析证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("AppendCertsFromPEM error ")
		return
	}
	c := credentials.NewTLS(&tls.Config{
		//设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cred},
		//要求必须校验客户端的证书
		ClientAuth: tls.RequireAndVerifyClientCert,
		//设置根证书的集合，校验方式使用ClientAuth设定的模式
		ClientCAs: certPool,
	})

	s := grpc.NewServer(grpc.Creds(c))
	tls_demo.RegisterHelloServiceServer(s, &TLSService{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("close err:%v", err)
		return
	}

}
