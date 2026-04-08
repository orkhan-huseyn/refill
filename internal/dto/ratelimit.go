package dto

import "time"

type RateLimitResponse struct {
	Allowed   bool      `json:"allowed"`
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset_time"`
}
