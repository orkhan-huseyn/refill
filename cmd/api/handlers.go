package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type RateLimitResponse struct {
	Allowed   bool      `json:"allowed"`
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset_time"`
}

var (
	capacity   = 10.0
	tokens     = capacity
	refillRate = 5.0
	lastRefill = time.Now()
)

func isAllowed(w http.ResponseWriter, r *http.Request) {
	refill()
	res := &RateLimitResponse{}

	if tokens > 1.0 {
		tokens -= 1.0
		res.Allowed = true
		w.WriteHeader(http.StatusOK)
	} else {
		res.Allowed = false
		w.WriteHeader(http.StatusTooManyRequests)
	}

	res.Remaining = int(tokens)
	res.ResetTime = getResetTime()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func refill() {
	now := time.Now()
	elapsed := now.Sub(lastRefill).Seconds()
	tokens = tokens + (elapsed * refillRate)
	if tokens > capacity {
		tokens = capacity
	}
	lastRefill = now
}

func getResetTime() time.Time {
	missingTokens := capacity - tokens
	secondsUntilFull := float64(missingTokens) / refillRate
	return time.Now().Add(time.Duration(secondsUntilFull * float64(time.Second)))
}
