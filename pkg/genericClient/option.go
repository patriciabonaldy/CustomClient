package genericClient

import (
	"math"
	"time"
)

type Options struct {
	TimeDuration     int
	MaxRetryCount    int
	ShouldRetry      bool
	CalculateBackoff func(attemptCount int) time.Duration
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

func WithBackoffPolicy() Option {
	min := 1
	return func(r *Options) {
		r.CalculateBackoff = func(attemptCount int) time.Duration {
			nextWait := time.Duration(math.Pow(2, float64(min)))
			min++

			return nextWait
		}
	}
}
