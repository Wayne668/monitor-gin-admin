package util

import (
	"context"
	"time"
)

type TokenBucket struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewTokenBucket(rate time.Duration, capacity int) *TokenBucket {
	ctx, cancel := context.WithCancel(context.Background())
	tb := &TokenBucket{
		tokens:   make(chan struct{}, capacity),
		rate:     rate,
		capacity: capacity,
		ctx:      ctx,
		cancel:   cancel,
	}

	for i := 0; i < capacity; i++ {
		tb.tokens <- struct{}{}
	}

	go tb.refill()
	return tb
}

func (tb *TokenBucket) refill() {
	ticker := time.NewTicker(tb.rate)
	defer ticker.Stop()

	for {
		select {
		case <-tb.ctx.Done():
			return
		case <-ticker.C:
			select {
			case tb.tokens <- struct{}{}:
			default:
			}
		}
	}
}

func (tb *TokenBucket) Take() {
	<-tb.tokens
}

func (tb *TokenBucket) TakeWithContext(ctx context.Context) error {
	select {
	case <-tb.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (tb *TokenBucket) Stop() {
	tb.cancel()
}
