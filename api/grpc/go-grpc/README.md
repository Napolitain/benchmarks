# Go gRPC Benchmark

A simple gRPC "Hello World" server for benchmarking.

## Run

```bash
go run main.go
```

Server listens on port 50051 (single-threaded via GOMAXPROCS=1).

## Generate proto (requires protoc)

```bash
protoc --go_out=. --go-grpc_out=. proto/helloworld.proto
```

## Test with grpcurl

```bash
grpcurl -plaintext -d '{"name":"World"}' localhost:50051 helloworld.Greeter/SayHello
```
