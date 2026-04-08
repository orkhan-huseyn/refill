package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/orkhan-huseyn/refill/internal/dto"
	"github.com/orkhan-huseyn/refill/internal/limiter"
)

var rateLimiter = limiter.New(10.0, 0.5)

func isAllowed(w http.ResponseWriter, r *http.Request) {
	allowed := rateLimiter.Allow()
	resetTime := rateLimiter.ResetTime()
	remaining := rateLimiter.Remaining()

	res := &dto.RateLimitResponse{
		Allowed:   allowed,
		ResetTime: resetTime,
		Remaining: remaining,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(rateLimiter.Capacity())))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(resetTime.UnixMilli())))
	w.Header().Set("Retry-After", strconv.Itoa(int(rateLimiter.RetryAfter().Seconds())))

	if allowed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
	json.NewEncoder(w).Encode(res)
}
