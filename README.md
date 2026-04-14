# Refill

Distributed rate limiter written in Go (for learning purposes).

### Running locally

Run with in memory storage: `go run ./...`
Run with redis storage:
```bash
$ docker compose up -d
$ go run ./... -storage=redis -redis.url=redis://default:my_password_here@localhost:6379/1
```

### Generate protobuf

```bash
$ brew install bufbuild/buf/buf
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$ export PATH="$PATH:$(go env GOPATH)/bin" # if you haven't done so
$ buf generate
```
