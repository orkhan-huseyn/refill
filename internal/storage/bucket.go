package storage

import (
	"time"
)

type Bucket struct {
	capacity   float64
	tokens     float64
	refillRate float64
	lastRefill time.Time
}

func NewBucket(capacity, rate float64) *Bucket {
	return &Bucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

func (b *Bucket) Refill() {
	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()

	b.tokens = min(b.capacity, b.tokens+(elapsed*b.refillRate))
	b.lastRefill = now
}

func (b *Bucket) RetryAfter() time.Duration {
	if b.tokens >= 1.0 {
		return 0
	}

	needed := 1.0 - b.tokens
	seconds := needed / b.refillRate
	return time.Duration(seconds * float64(time.Second))
}

func (b *Bucket) ResetTime() time.Time {
	missingTokens := b.capacity - b.tokens
	secondsUntilFull := float64(missingTokens) / b.refillRate
	return time.Now().Add(time.Duration(secondsUntilFull * float64(time.Second)))
}

func (b *Bucket) Consume(cost float64) {
	b.tokens -= cost
}

func (b *Bucket) Capacity() float64 {
	return b.capacity
}

func (b *Bucket) Remaining() float64 {
	return b.tokens
}
