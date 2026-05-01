package limiter

import (
	"context"

	"github.com/orkhan-huseyn/refill/config"
	"github.com/orkhan-huseyn/refill/internal/counter"
	"github.com/orkhan-huseyn/refill/internal/ruleprovider"
)

type Limiter struct {
	counter      counter.RateLimitCounter
	ruleProvider ruleprovider.RuleProvider
}

func NewLimiter(cfg config.Config) *Limiter {
	// TODO: move it to factory method and handle errors (e.g. redisurl is not passed)
	var counterInstance counter.RateLimitCounter
	switch cfg.RateLimit.Type {
	case config.RateLimitLocal:
		counterInstance = counter.NewInMemoryCounter()
	case config.RateLimitGlobal:
		counterInstance = counter.NewRedisCounter(cfg.RateLimit.Redis)
	}

	// TODO: move it to factory method and handle errors
	var ruleProviderInstance ruleprovider.RuleProvider
	switch cfg.RuleProvider.Type {
	case config.ProviderTypeStatic:
		ruleProviderInstance = ruleprovider.NewStaticProvider(cfg)
	}

	// TODO: is this right place to do this? (also handle error)
	go ruleProviderInstance.PopulateCache()

	return &Limiter{
		counterInstance,
		ruleProviderInstance,
	}
}

func (l *Limiter) Allow(ctx context.Context, key, namespace string, cost int) (counter.RateLimitResult, error) {
	compositeKey := key + ":" + namespace
	rule, err := l.ruleProvider.GetRule(namespace)
	if err != nil {
		return counter.RateLimitResult{}, err
	}
	return l.counter.Take(ctx, compositeKey, cost, rule.Burst, rule.Rate)
}
