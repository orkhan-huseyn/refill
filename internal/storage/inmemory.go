package storage

import (
	"context"
	"runtime"

	"github.com/orkhan-huseyn/refill/internal/shardedmap"
)

type InMemoryStore struct {
	cache shardedmap.ShardedMap[*Bucket]
}

func NewInMemoryStore() InMemoryStore {
	// TODO: container aware?
	shardCount := runtime.NumCPU()
	return InMemoryStore{
		cache: shardedmap.New[*Bucket](shardCount),
	}
}

func (s InMemoryStore) Take(ctx context.Context, key string, amount int) (RateLimitResult, error) {
	var res RateLimitResult
	if err := ctx.Err(); err != nil {
		return res, err
	}

	bucket, ok := s.cache.Get(key)
	if !ok {
		// TODO: where to get those arguments? hardcoded doesn't look good
		bucket = NewBucket(5.0, 0.5)
		s.cache.Put(key, bucket)
	}

	bucket.Refill()
	res = RateLimitResult{
		Limit:      int(bucket.capacity),
		Remaining:  int(bucket.tokens),
		RetryAfter: bucket.RetryAfter(float64(amount)),
		ResetTime:  bucket.ResetTime(),
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
