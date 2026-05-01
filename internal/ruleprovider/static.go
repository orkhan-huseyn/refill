package ruleprovider

import (
	"fmt"
	"runtime"

	"github.com/orkhan-huseyn/refill/config"
	"github.com/orkhan-huseyn/refill/internal/dto"
	"github.com/orkhan-huseyn/refill/internal/shardedmap"
)

type StaticRuleProvider struct {
	cache shardedmap.ShardedMap[*dto.RateLimitRule]
	rules []dto.RateLimitRule
}

func NewStaticProvider(cfg config.Config) StaticRuleProvider {
	// TODO: container aware?
	shardCount := runtime.NumCPU()

	return StaticRuleProvider{
		cache: shardedmap.New[*dto.RateLimitRule](shardCount),
		rules: cfg.RuleProvider.Rules,
	}
}

func (e StaticRuleProvider) PopulateCache() error {
	for _, rule := range e.rules {
		e.cache.Put(rule.Namespace, &rule)
	}
	return nil
}

func (e StaticRuleProvider) GetRule(namespace string) (dto.RateLimitRule, error) {
	rule, exists := e.cache.Get(namespace)
	if !exists {
		return dto.RateLimitRule{}, fmt.Errorf("no rule exists for namaspace '%s'", namespace)
	}
	return *rule, nil
}
