package enforcer

type RuleEnforcer interface {
	GetRule(namespace string) (RateLimitRule, error)
}

type RateLimitRule struct {
	Namespace string  `yaml:"namespace"`
	Burst     float64 `yaml:"burst"`
	Rate      float64 `yaml:"rate"`
}
