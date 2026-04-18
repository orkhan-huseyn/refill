package enforcer

import (
	"fmt"
	"runtime"

	"github.com/orkhan-huseyn/refill/config"
	"github.com/orkhan-huseyn/refill/internal/dto"
	"github.com/orkhan-huseyn/refill/internal/shardedmap"
)

type StaticRuleEnforcer struct {
	cache shardedmap.ShardedMap[*dto.RateLimitRule]
	rules []dto.RateLimitRule
}

func NewStaticEnforcer(cfg config.Config) StaticRuleEnforcer {
	// TODO: container aware?
	shardCount := runtime.NumCPU()

	return StaticRuleEnforcer{
		cache: shardedmap.New[*dto.RateLimitRule](shardCount),
		rules: cfg.Enforcer.Rules,
	}
}

func (e StaticRuleEnforcer) PopulateCache() error {
	for _, rule := range e.rules {
		e.cache.Put(rule.Namespace, &rule)
	}
	return nil
}

func (e StaticRuleEnforcer) GetRule(namespace string) (dto.RateLimitRule, error) {
	rule, exists := e.cache.Get(namespace)
	if !exists {
		return dto.RateLimitRule{}, fmt.Errorf("no rule exists for namaspace '%s'", namespace)
	}
	return *rule, nil
}
