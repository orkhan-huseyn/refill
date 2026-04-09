package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/orkhan-huseyn/refill/internal/dto"
	"github.com/orkhan-huseyn/refill/internal/limiter"
)

var rateLimiter = limiter.NewLimiter()

func isAllowed(w http.ResponseWriter, r *http.Request) {
	req := &dto.RateLimitRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := rateLimiter.Allow(req.Key, req.Namespace, req.Cost)

	res := &dto.RateLimitResponse{
		Allowed:   result.Allowed,
		ResetTime: result.Reset,
		Remaining: int(result.Remaining),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(result.Capacity)))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(result.Remaining)))
	w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(result.Reset.UnixMilli())))
	w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))

	if result.Allowed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
	json.NewEncoder(w).Encode(res)
}
