package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	mu      sync.RWMutex
	storage map[string]*Bucket
}

type Result struct {
	Allowed    bool
	Remaining  float64
	Capacity   float64
	Reset      time.Time
	RetryAfter time.Duration
}

func NewLimiter() *Limiter {
	return &Limiter{
		storage: make(map[string]*Bucket),
	}
}

func (l *Limiter) Allow(key, namespace string, cost int) *Result {
	b := l.getBucket(key, namespace)
	b.refill()

	res := &Result{}
	res.Capacity = b.capacity
	res.Remaining = b.tokens
	res.RetryAfter = b.RetryAfter()
	res.Reset = b.ResetTime()

	// TODO: what is cost is zero?
	if b.tokens >= 1.0 {
		b.tokens -= float64(cost)
		res.Allowed = true
		res.Remaining = b.tokens
		return res
	}

	res.Allowed = false
	return res
}

func (l *Limiter) getBucket(key, namespace string) *Bucket {
	compositeKey := key + ":" + namespace

	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.storage[compositeKey]
	if !ok {
		b = NewBucket(5.0, 0.5)
		l.storage[compositeKey] = b
	}

	return b
}
