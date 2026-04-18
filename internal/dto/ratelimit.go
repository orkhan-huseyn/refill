package dto

// TODO: should this be located here
type RateLimitRule struct {
	Namespace string  `yaml:"namespace"`
	Burst     float64 `yaml:"burst"`
	Rate      float64 `yaml:"rate"`
}
