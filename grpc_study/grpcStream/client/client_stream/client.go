package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpcStream/pb/stream"
	"io"
	"log"
	"sync"
	"time"
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

func bothStreamDemo(client pb.StreamServiceClient) {
	var wg sync.WaitGroup
	// 2. 调用方法获取stream
	stream, err := client.BothStream(context.Background())
	if err != nil {
		panic(err)
	}
	// 3.开两个goroutine 分别用于Recv()和Send()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("服务端已关闭")
				break
			}
			if err != nil {
				continue
			}
			fmt.Printf("客户端收到服务端的数据 :%v \n", req.GetName())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 暂且设置只给服务端流式通信3次
		for i := 0; i < 3; i++ {
			err := stream.Send(&pb.StreamRequest{Name: "哟，服务端的小刁毛"})
			if err != nil {
				log.Printf("发送失败:%v\n", err)
			}
			time.Sleep(time.Second)
		}
		// 4. 发送完毕关闭stream
		err := stream.CloseSend()
		if err != nil {
			log.Printf("客户端发送错误:%v\n", err)
			return
		}
	}()
	wg.Wait()
}

func main() {
	conn, err := grpc.Dial(Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)
	serverStreamDemo(client, "舔狗，过来舔我！")
	//clientStreamDemo(client)
	//bothStreamDemo(client)
}
