package main

import (
	"fmt"
	"google.golang.org/grpc"
	sp "grpcStream/pb/stream"
	"io"
	"log"
	"net"
)

const Port = ":8999"

type StreamSer struct {
	sp.UnimplementedStreamServiceServer
}

// ServerStream 服务端流模式
func (s *StreamSer) ServerStream(in *sp.StreamRequest, out sp.StreamService_ServerStreamServer) error {
	log.Printf("收到客户端的请求 %v", in.GetName())
	// 返回多份数据给client，假设我们模拟返回10条数据给client
	for i := 0; i < 10; i++ {
		err := out.Send(&sp.StreamResponse{Name: fmt.Sprintf("丢雷- %d", i)})
		if err != nil {
			log.Fatalf("Server Stream Send error:%v", err)
			return err
		}
	}
	// 返回nil表示完成响应
	return nil
}

//ClientStream：客户端流式 RPC
func (s *StreamSer) ClientStream(clientStr sp.StreamService_ClientStreamServer) error {
	for {
		r, err := clientStr.Recv()
		if err == io.EOF {
			// SendAndClose 返回并关闭连接
			// 在客户端发送完毕后服务端即可返回响应
			return clientStr.SendAndClose(&sp.StreamResponse{Name: "客户端就是个纯纯的舔狗"})
		}
		if err != nil {
			return err
		}
		log.Printf("收到客户端的流: %s", r.Name)
	}
	return nil
}

//BothStream：双向流式 RPC
func (s *StreamSer) BothStream(in sp.StreamService_BothStreamServer) error {
	return nil
}

func main() {
	server := grpc.NewServer()                           // 新建 grpc server 对象
	sp.RegisterStreamServiceServer(server, &StreamSer{}) // 注册

	listen, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("grpc.Dial err :%v", err)
	}
	defer listen.Close()

	server.Serve(listen)

}
