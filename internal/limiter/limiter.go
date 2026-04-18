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
	case config.RateLimitLocal:
		storageToUse = storage.NewInMemoryStore()
	case config.RateLimitGlobal:
		storageToUse = storage.NewRedisStore(cfg.RateLimit.Redis)
	}

	// TODO: move it to factory method and handle errors
	var enforcerToUse enforcer.RuleEnforcer
	switch cfg.Enforcer.Type {
	case config.TypeStatic:
		enforcerToUse = enforcer.NewStaticEnforcer(cfg)
	}

	// TODO: is this right place to do this? (also handle error)
	go enforcerToUse.PopulateCache()

	return &Limiter{
		storage:  storageToUse,
		enforcer: enforcerToUse,
	}
}

func (l *Limiter) Allow(ctx context.Context, key, namespace string, cost int) (storage.RateLimitResult, error) {
	compositeKey := key + ":" + namespace
	rule, err := l.enforcer.GetRule(namespace)
	if err != nil {
		return storage.RateLimitResult{}, err
	}
	return l.storage.Take(ctx, compositeKey, cost, rule.Burst, rule.Rate)
}
