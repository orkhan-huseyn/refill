package config

import "github.com/orkhan-huseyn/refill/internal/dto"

type RateLimitType string

const (
	RateLimitLocal  RateLimitType = "local"
	RateLimitGlobal RateLimitType = "global"
)

// TODO: we'll use redis for enforcer as well, so maybe some global config
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RateLimitConfig struct {
	Type  RateLimitType `yaml:"type"`
	Redis RedisConfig   `yaml:"redis,omitempty"`
}

type EnforcerType string

const (
	TypeStatic   EnforcerType = "static"
	TypePostgres EnforcerType = "postgres"
)

type EnforcerConfig struct {
	Type   EnforcerType        `yaml:"type"`
	Rules  []dto.RateLimitRule `yaml:"rules,omitempty"`
	DBConn string              `yaml:"dbconn,omitempty"` // TODO: make it PostgresConfig struct
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	RateLimit RateLimitConfig `yaml:"ratelimit"`
	Enforcer  EnforcerConfig  `yaml:"enforcer"`
}
