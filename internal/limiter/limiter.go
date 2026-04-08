package limiter

import (
	"math"
	"time"
)

type Limiter struct {
	capacity   float64
	tokens     float64
	refillRate float64
	lastRefill time.Time
}

func New(capacity, rate float64) *Limiter {
	return &Limiter{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	l.refill()

	if l.tokens >= 1.0 {
		l.tokens--
		return true
	}
	return false
}

func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()

	l.tokens = math.Min(l.capacity, l.tokens+(elapsed*l.refillRate))
	l.lastRefill = now
}

func (l *Limiter) RetryAfter() time.Duration {
	if l.tokens >= 1.0 {
		return 0
	}

	needed := 1.0 - l.tokens
	seconds := needed / l.refillRate
	return time.Duration(seconds * float64(time.Second))
}

func (l *Limiter) ResetTime() time.Time {
	missingTokens := l.capacity - l.tokens
	secondsUntilFull := float64(missingTokens) / l.refillRate
	return time.Now().Add(time.Duration(secondsUntilFull * float64(time.Second)))
}

func (l *Limiter) Capacity() int {
	return int(l.capacity)
}

func (l *Limiter) Remaining() int {
	return int(l.tokens)
}
