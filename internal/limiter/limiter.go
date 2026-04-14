package limiter

import (
	"context"

	"github.com/orkhan-huseyn/refill/internal/storage"
)

type Limiter struct {
	storage storage.RateLimitStore
}

func NewLimiter(storageType, redisUrl string) *Limiter {
	// TODO: move it to factory method and handle errors (e.g. redisurl is not passed)
	var storageToUse storage.RateLimitStore
	switch storageType {
	case "inmemory":
		storageToUse = storage.NewInMemoryStore()
	case "redis":
		storageToUse = storage.NewRedisStore(redisUrl)
	}
	return &Limiter{
		storage: storageToUse,
	}
}

// TODO: fetch limit and rate from rule storage
var limit = 5.0
var rate = 0.5

func (l *Limiter) Allow(ctx context.Context, key, namespace string, cost int) (storage.RateLimitResult, error) {
	compositeKey := key + ":" + namespace
	return l.storage.Take(ctx, compositeKey, cost, limit, rate)
}
