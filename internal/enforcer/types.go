package enforcer

import "github.com/orkhan-huseyn/refill/internal/dto"

type RuleEnforcer interface {
	PopulateCache() error
	GetRule(namespace string) (dto.RateLimitRule, error)
}
