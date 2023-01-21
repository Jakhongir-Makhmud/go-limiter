package golimiter

import (
	"context"
	"go-limiter/stateStorage"
	"time"
)

type StateStorage interface {
	GetState(ctx context.Context, key string) (uint64, error)
	SetState(ctx context.Context, key string, val uint64, dur time.Duration) error
}

type Config struct {
	RequestsPerDuration    uint64
	LimitDuration          time.Duration
	ExcludedRequestMethods []string
	StateStorage           StateStorage
}

type RateLimiter struct {
	requestPerDuration     uint64
	LimitDuration          time.Duration
	excludedRequestMethods []string
	StateStorage           StateStorage
}

func New(cfg Config) *RateLimiter {

	var stateStr = cfg.StateStorage

	if stateStr == nil {
		stateStr = stateStorage.New()
	}

	return &RateLimiter{
		requestPerDuration:     cfg.RequestsPerDuration,
		LimitDuration:          cfg.LimitDuration,
		excludedRequestMethods: cfg.ExcludedRequestMethods,
		StateStorage:           stateStr,
	}
}

func (rl *RateLimiter) AddOne(key string) error {

	rate, err := rl.StateStorage.GetState(context.Background(), key)
	if err != nil {
		return err
	}

	return rl.StateStorage.SetState(context.Background(), key, rate+1, rl.LimitDuration)
}

func (rl *RateLimiter) IsAccepted(key string) bool {
	rate, err := rl.StateStorage.GetState(context.Background(), key)
	if err != nil {
		return false
	}

	if rl.requestPerDuration >= rate {
		return true
	}

	return false
}
