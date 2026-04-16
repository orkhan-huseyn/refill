package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/orkhan-huseyn/refill/config"
	pb "github.com/orkhan-huseyn/refill/gen/go/v1"
	ratelimitsrv "github.com/orkhan-huseyn/refill/internal/server"
)

var cfg config.Config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/refill/")
	viper.AddConfigPath("$HOME/.refill")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into Config struct: %v", err)
	}
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", cfg.Server.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	serviceImpl := ratelimitsrv.NewRateLimitServer(cfg)
	pb.RegisterRateLimitServiceServer(server, serviceImpl)

	go func() {
		log.Printf("Starting server on port :%s", cfg.Server.Addr)
		if err := server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatal(err)
		}
	}()

	<-shutdown
	server.GracefulStop()
	log.Println("Server stopped")
}
