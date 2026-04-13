package storage

import (
	"context"
	"time"
)

type RateLimitStore interface {
	Take(ctx context.Context, key string, amount int, limit float64, rate float64) (RateLimitResult, error)
}

type RateLimitResult struct {
	Allowed    bool
	Limit      int
	Remaining  int
	RetryAfter time.Duration
	ResetTime  time.Time
}
