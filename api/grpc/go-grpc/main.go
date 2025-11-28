package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"

	pb "go-grpc/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, World!"}, nil
}

func main() {
	runtime.GOMAXPROCS(1)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	fmt.Println("gRPC server listening on :50051 (single-threaded)")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
