package counter

import (
	"context"
	"runtime"

	"github.com/orkhan-huseyn/refill/internal/shardedmap"
)

type InMemoryCounter struct {
	cache shardedmap.ShardedMap[*Bucket]
}

func NewInMemoryCounter() InMemoryCounter {
	// TODO: container aware?
	shardCount := runtime.NumCPU()
	return InMemoryCounter{
		cache: shardedmap.New[*Bucket](shardCount),
	}
}

func (c InMemoryCounter) Take(ctx context.Context, key string, amount int, limit float64, rate float64) (RateLimitResult, error) {
	var res RateLimitResult
	if err := ctx.Err(); err != nil {
		return res, err
	}

	bucket, ok := c.cache.Get(key)
	if !ok {
		bucket = NewBucket(limit, rate)
		c.cache.Put(key, bucket)
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
