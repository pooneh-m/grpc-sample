package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	pb "grpc-sample/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":12343"
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

	opts := getServerOptions()
	s := grpc.NewServer(opts)

	pb.RegisterHelloWorldServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getServerOptions() grpc.ServerOption {
	tlsCer, err := tls.LoadX509KeyPair("../creds/service.pem", "../creds/service.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials: %v", err)
	}

	certPool := x509.NewCertPool()
	clientCA, err := ioutil.ReadFile("../creds/client.pem")
	if err != nil {
		log.Fatalf("failed to read client ca cert: %s", err)
	}
	ok := certPool.AppendCertsFromPEM(clientCA)
	if !ok {
		log.Fatal("failed to append client certs")
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{tlsCer},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		// GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
		// 	return &tls.Config{
		// 		Certificates: []tls.Certificate{tlsCer},
		// 		ClientAuth:   tls.RequireAndVerifyClientCert,
		// 		ClientCAs:    certPool,
		// 	}, nil
		// },
	}
	// Add options for creds and OpenCensus stats handler to enable stats and tracing.
	return grpc.Creds(credentials.NewTLS(cfg))
}
