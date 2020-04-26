In this repo the servers are in golang and clients are in C#.

To run servers, go to server directory and run:
```
go run main.go
```
To run the clients, go to client directory and run:
```
dotnet run
```
With hellomtls sample the request fails.

Generate golang proto
```
protoc ../sample.proto --proto_path=.. --go_out=plugins=grpc:.
```