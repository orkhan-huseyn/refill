package server

import (
	"context"

	pb "github.com/orkhan-huseyn/refill/gen/go/v1"
	"github.com/orkhan-huseyn/refill/internal/limiter"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RateLimitServer struct {
	pb.UnimplementedRateLimitServiceServer
	limiter *limiter.Limiter
}

func NewRateLimitServer() *RateLimitServer {
	return &RateLimitServer{
		limiter: limiter.NewLimiter(),
	}
}

func (s *RateLimitServer) IsAllowed(ctx context.Context, req *pb.RateLimitRequest) (*pb.RateLimitResponse, error) {
	res, err := s.limiter.Allow(ctx, req.Key, req.Namespace, int(req.Cost))
	if err != nil {
		return nil, err
	}

	return &pb.RateLimitResponse{
		Allowed:    res.Allowed,
		Remaining:  int32(res.Remaining),
		ResetTime:  timestamppb.New(res.ResetTime),
		RetryAfter: durationpb.New(res.RetryAfter),
	}, nil
}
