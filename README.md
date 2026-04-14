# Refill

Distributed rate limiter written in Go (for learning purposes).

### Generate protobuf

```bash
$ brew install bufbuild/buf/buf
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$ export PATH="$PATH:$(go env GOPATH)/bin" # if you haven't done so
$ buf generate
```
