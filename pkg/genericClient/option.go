package genericClient

import (
	"time"
)

type Options struct {
	TimeDuration     int
	MaxRetryCount    int
	ShouldRetry      bool
	CalculateBackoff time.Duration
}

type Option func(*Options)

func WithTimeDuration(timeDuration int) Option {
	if timeDuration < 0 {
		timeDuration = 10
	}
	return func(r *Options) {
		r.TimeDuration = timeDuration
	}
}

func WithMaxRetryCount(maxRetryCount int) Option {
	if maxRetryCount < 0 {
		maxRetryCount = 0
	}
	return func(r *Options) {
		r.MaxRetryCount = maxRetryCount
	}
}

func WithRetryPolicy(retryPolicy bool) Option {
	return func(r *Options) {
		r.ShouldRetry = retryPolicy
	}
}

func WithBackoffPolicy(backoffPolicy time.Duration) Option {
	return func(roundtripper *Options) {
		roundtripper.CalculateBackoff = backoffPolicy
	}
}
