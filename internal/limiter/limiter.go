package limiter

import (
	"context"

	"github.com/orkhan-huseyn/refill/internal/storage"
)

type Limiter struct {
	storage storage.RateLimitStore
}

func NewLimiter() *Limiter {
	return &Limiter{
		storage: storage.NewRedisStore("redis://default:my_password_here@localhost:6379/1"),
	}
}

// TODO: fetch limit and rate from rule storage
var limit = 5.0
var rate = 0.5

func (l *Limiter) Allow(ctx context.Context, key, namespace string, cost int) (storage.RateLimitResult, error) {
	compositeKey := key + ":" + namespace
	return l.storage.Take(ctx, compositeKey, cost, limit, rate)
}
