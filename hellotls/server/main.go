package main

import (
	"context"
	"log"
	"net"

	pb "grpc-sample/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	creds, err := credentials.NewServerTLSFromFile("../creds/service.pem", "../creds/service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterHelloWorldServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/*func getServerOptions() ([]grpc.ServerOption, error) {
	tlsCer, err := tls.LoadX509KeyPair(tlsDir+"tls.crt", tlsDir+"tls.key")
	if err != nil {
		logger.WithError(err).Fatal("failed to generate credentials")
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{tlsCer},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
			return &tls.Config{
				Certificates: []tls.Certificate{tlsCer},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    h.caCertPool,
			}, nil
		},
	}
	// Add options for creds and OpenCensus stats handler to enable stats and tracing.
	return []grpc.ServerOption{grpc.Creds(credentials.NewTLS(cfg)), grpc.StatsHandler(&ocgrpc.ServerHandler{})}
}*/
