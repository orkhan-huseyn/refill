package limiter

import (
	"context"

	"github.com/orkhan-huseyn/refill/config"
	"github.com/orkhan-huseyn/refill/internal/enforcer"
	"github.com/orkhan-huseyn/refill/internal/storage"
)

type Limiter struct {
	storage  storage.RateLimitStore
	enforcer enforcer.RuleEnforcer
}

func NewLimiter(cfg config.Config) *Limiter {
	// TODO: move it to factory method and handle errors (e.g. redisurl is not passed)
	var storageToUse storage.RateLimitStore
	switch cfg.RateLimit.Type {
	case "local":
		storageToUse = storage.NewInMemoryStore()
	case "global":
		storageToUse = storage.NewRedisStore(cfg.RateLimit.Redis)
	}
	return &Limiter{
		storage: storageToUse,
	}
}

// TODO: fetch limit and rate from rule storage
var burst = 5.0
var rate = 0.5

func (l *Limiter) Allow(ctx context.Context, key, namespace string, cost int) (storage.RateLimitResult, error) {
	compositeKey := key + ":" + namespace
	return l.storage.Take(ctx, compositeKey, cost, burst, rate)
}
