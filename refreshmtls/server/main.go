package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"time"

	pb "grpc-sample/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":12344"
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

	tlsConfig := getTLSConfig()
	opts := grpc.Creds(credentials.NewTLS(tlsConfig))
	s := grpc.NewServer(opts)

	pb.RegisterHelloWorldServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getTLSConfig() *tls.Config {
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

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCer},
		ClientAuth:   tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			log.Print("Verifying peer certificate")
			opts := x509.VerifyOptions{
				Roots:         certPool,
				CurrentTime:   time.Now(),
				Intermediates: x509.NewCertPool(),
				KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}

			for _, cert := range rawCerts[1:] {
				opts.Intermediates.AppendCertsFromPEM(cert)
			}

			c, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return errors.New("tls: failed to verify client certificate: " + err.Error())
			}
			_, err = c.Verify(opts)
			if err != nil {
				return errors.New("tls: failed to verify client certificate: " + err.Error())
			}
			return nil
		},
	}
}
