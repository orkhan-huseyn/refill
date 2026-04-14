package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/orkhan-huseyn/refill/gen/go/v1"
	ratelimitsrv "github.com/orkhan-huseyn/refill/internal/server"
)

var port = flag.Int("port", 50051, "The server port")
var storage = flag.String("storage", "inmemory", "Storage to use, (redis|memcached|inmemory)")
var redisUrl = flag.String("redis.url", "", "Redis url e.g. redis://user:pass@localhost:6379/<db>")

func main() {
	flag.Parse()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	serviceImpl := ratelimitsrv.NewRateLimitServer(*storage, *redisUrl)
	pb.RegisterRateLimitServiceServer(server, serviceImpl)

	go func() {
		log.Printf("Starting server on port :%d", *port)
		if err := server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatal(err)
		}
	}()

	<-shutdown
	server.GracefulStop()
	log.Println("Server stopped")
}
