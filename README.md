In this repo the servers are in golang and clients are in C#.

To run servers, go to server directory and run:
```
go run main.go
```
To run the clients, go to client directory and run:
```
dotnet run
```

Generate golang proto
```
protoc ../sample.proto --proto_path=.. --go_out=plugins=grpc:.
```

Certificates are generated using:
```
TLS_KEY_FILE=service.key
TLS_CERT_FILE=service.pem

openssl req -nodes -new -newkey rsa:2048 \
    -keyout ${TLS_KEY_FILE} \
    -out tls.csr \
    -subj "/CN=localhost"

openssl x509 -req -days 365 -in tls.csr \
    -signkey ${TLS_KEY_FILE} \
    -out ${TLS_CERT_FILE}

KEY_FILE=client.key
CERT_FILE=client.pem

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE}
```