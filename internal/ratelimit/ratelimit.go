// Package ratelimit provides token-bucket rate limiting for LLM requests.
package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// Limiter wraps a golang.org/x/time/rate.Limiter and provides
// cancellation-aware waiting for LLM requests.
type Limiter struct {
	lim *rate.Limiter
}

// New creates a Limiter with the given requests-per-second and burst.
// If rate is zero or burst is zero, the limiter allows all requests (no-op).
func New(rateLimit, burst int) *Limiter {
	if rateLimit == 0 || burst == 0 {
		return &Limiter{lim: nil}
	}
	return &Limiter{
		lim: rate.NewLimiter(rate.Every(time.Duration(rateLimit)*time.Second), burst),
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns nil if the limiter is disabled (nil lim).
func (l *Limiter) Wait(ctx context.Context) error {
	if l.lim == nil {
		return nil
	}
	return l.lim.Wait(ctx)
}

// Allow reports whether a token is currently available without blocking.
func (l *Limiter) Allow() bool {
	if l.lim == nil {
		return true
	}
	return l.lim.Allow()
}

// Reserve reserves a token and returns a Reservation that can be
// cancelled or waited on. Returns nil if the limiter is disabled.
func (l *Limiter) Reserve() *rate.Reservation {
	if l.lim == nil {
		return nil
	}
	return l.lim.Reserve()
}

// SetLimit updates the rate (requests per second) dynamically.
func (l *Limiter) SetLimit(rateLimit int) {
	if l.lim == nil {
		return
	}
	l.lim.SetLimit(rate.Every(time.Duration(rateLimit) * time.Second))
}

// SetBurst updates the burst size dynamically.
func (l *Limiter) SetBurst(burst int) {
	if l.lim == nil {
		return
	}
	l.lim.SetBurst(burst)
}