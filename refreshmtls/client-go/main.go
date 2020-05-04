package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	pb "grpc-sample/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	keyFile := "../creds/client.key"
	certFile := "../creds/client.pem"
	cacertFile := "../creds/service.pem"

	endpoint := "localhost:12344"
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		panic(err)
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	cacert, err := ioutil.ReadFile(cacertFile)
	if err != nil {
		panic(err)
	}

	dialOpts, err := createRemoteClusterDialOption(cert, key, cacert)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(endpoint, dialOpts)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	grpcClient := pb.NewHelloWorldServiceClient(conn)
	response, err := grpcClient.Test(context.Background(), &pb.TestRequest{Query: "go client mTLS"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("response: %s\n", response.String())
}

// createRemoteClusterDialOption creates a grpc client dial option with TLS configuration.
func createRemoteClusterDialOption(clientCert, clientKey, caCert []byte) (grpc.DialOption, error) {
	// Load client cert
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	if len(caCert) != 0 {
		// Load CA cert, if provided and trust the server certificate.
		// This is required for self-signed certs.
		tlsConfig.RootCAs = x509.NewCertPool()
		if !tlsConfig.RootCAs.AppendCertsFromPEM(caCert) {
			return nil, errors.New("only PEM format is accepted for server CA")
		}
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
}
