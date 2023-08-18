package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpcStream/pb/stream"
	"io"
	"log"
)

const Port = ":8999"

func serverStreamDemo(client pb.StreamServiceClient, name string) error {
	stream, err := client.ServerStream(context.Background(), &pb.StreamRequest{Name: name})
	if err != nil {
		log.Fatalf("调用服务端stream失败:%v", err)
		return err
	}
	for {
		resp, err := stream.Recv()
		// 4. err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			log.Println("server closed")
			break
		}
		if err != nil {
			log.Printf("Recv error:%v", err)
			continue
		}
		log.Printf("Recv data:%v", resp.GetName())
	}
	return nil
}

// 客户端流模式
func clientStreamDemo(client pb.StreamServiceClient) error {
	// 循环接受客户端的信息
	stream, err := client.ClientStream(context.Background())
	if err != nil {
		return err
	}
	// 客户端发送10次，源源不断的给服务端请求
	for i := 0; i <= 10; i++ {
		// 通过 Send 方法不断推送数据到服务端
		err := stream.Send(&pb.StreamRequest{Name: fmt.Sprintf("客服端在第%d次请求服务端", i)})
		if err != nil {
			return err
		}
	}

	// 发送完成后通过stream.CloseAndRecv() 关闭stream并接收服务端返回结果
	// (服务端则根据err==io.EOF来判断client是否关闭stream)
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("服务端回应: %s", resp.GetName())
	return nil
}

func main() {
	conn, err := grpc.Dial(Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)
	//serverStreamDemo(client, "舔狗，过来舔我！")
	clientStreamDemo(client)

}
