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
)

const port = ":8999"

func main() {
	// 公钥中读取和解析公钥/私钥对
	cred, err := tls.LoadX509KeyPair("./conf/client.crt", "./conf/client.key")

	if err != nil {
		fmt.Println("LoadX509KeyPair error ", err)
		return
	}
	// 创建一组根证书
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		fmt.Println("ReadFile ca.crt error ", err)
		return
	}
	// 解析证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("certPool.AppendCertsFromPEM error ")
		return
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cred},
		ServerName:   "go-grpc-example",
		RootCAs:      certPool,
	})

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
