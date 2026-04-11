package storage

import (
	"context"
	"sync"
)

type InMemoryStore struct {
	mu    sync.RWMutex
	cache map[string]*Bucket
}

func NewInMemory() *InMemoryStore {
	return &InMemoryStore{
		mu:    sync.RWMutex{},
		cache: make(map[string]*Bucket),
	}
}

func (s *InMemoryStore) Take(ctx context.Context, key string, amount int) (RateLimitResult, error) {
	var res RateLimitResult
	if err := ctx.Err(); err != nil {
		return res, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, ok := s.cache[key]
	if !ok {
		// TODO: where to get those arguments? hardcoded doesn't look good
		bucket = NewBucket(5.0, 0.5)
		s.cache[key] = bucket
	}

	bucket.Refill()
	res = RateLimitResult{
		Limit:      int(bucket.capacity),
		Remaining:  int(bucket.tokens),
		RetryAfter: bucket.RetryAfter(),
	}

	// TODO: what if cost is zero?
	if bucket.tokens >= float64(amount) {
		bucket.tokens -= float64(amount)
		res.Allowed = true
		res.Remaining = int(bucket.tokens)
	} else {
		res.Remaining = int(bucket.tokens)
		res.Allowed = false
	}

	return res, nil
}
