package ruleprovider

import "github.com/orkhan-huseyn/refill/internal/dto"

type RuleProvider interface {
	PopulateCache() error
	GetRule(namespace string) (dto.RateLimitRule, error)
}
