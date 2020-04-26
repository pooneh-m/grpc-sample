package main

import (
	"context"
	"log"
	"net"

	pb "grpc-sample/sample"

	"google.golang.org/grpc"
)

const (
	port = ":12341"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedHelloWorldServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) Test(ctx context.Context, in *pb.TestRequest) (*pb.TestResponse, error) {
	log.Printf("Received: %v", in.GetQuery())
	return &pb.TestResponse{Message: "Got your query: " + in.GetQuery()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloWorldServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
