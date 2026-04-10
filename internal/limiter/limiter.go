package limiter

import (
	"time"

	"github.com/orkhan-huseyn/refill/internal/storage"
)

type Limiter struct {
	storage storage.Storage
}

type Result struct {
	Allowed    bool
	Remaining  int
	Capacity   int
	Reset      time.Time
	RetryAfter time.Duration
}

func NewLimiter() *Limiter {
	return &Limiter{
		storage: storage.NewInMemory(),
	}
}

func (l *Limiter) Allow(key, namespace string, cost int) *Result {
	b := l.storage.GetBucket(key, namespace)
	b.Refill()

	// TODO: maybe cleanup a bit?
	res := &Result{}
	res.Capacity = int(b.Capacity())
	res.Remaining = int(b.Remaining())
	res.RetryAfter = b.RetryAfter()
	res.Reset = b.ResetTime()

	// TODO: what is cost is zero?
	if float64(b.Remaining()) >= 1.0 {
		b.Consume(float64(cost))
		res.Allowed = true
		res.Remaining = int(b.Remaining())
		return res
	}

	res.Allowed = false
	return res
}
