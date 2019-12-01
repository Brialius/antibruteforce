package leakybucket

import (
	"context"
	"log"
	"sync"
	"time"
)

// Bucket struct
type Bucket struct {
	ID                  string
	rateLimit           uint64
	requests            uint64
	waitTicks           uint64
	waitTicksUntilPurge uint64
	done                chan struct{}
	mu                  sync.RWMutex
	duration            time.Duration
	inactiveCycles      uint64
}

// NewBucket Bucket constructor
func NewBucket(id string, rateLimit uint64, duration time.Duration, inactiveCycles uint64) *Bucket {
	b := &Bucket{ID: id, rateLimit: rateLimit, waitTicksUntilPurge: inactiveCycles * rateLimit,
		done: make(chan struct{}, 1), mu: sync.RWMutex{}, duration: duration, inactiveCycles: inactiveCycles}
	go func(bucket *Bucket) {
		log.Printf("Goroutine for ID: %s is created", b.ID)
		ticker := time.NewTicker(time.Duration(uint64(b.duration) / b.rateLimit))
		defer ticker.Stop()
		for {
			<-ticker.C
			bucket.mu.RLock()
			if bucket.waitTicks >= b.waitTicksUntilPurge {
				bucket.done <- struct{}{}
				bucket.mu.RUnlock()
				log.Printf("Exit from goroutine for ID: %s", b.ID)
				return
			}
			bucket.mu.RUnlock()
			bucket.mu.Lock()
			if bucket.requests > 0 {
				bucket.requests--
				bucket.waitTicks = 0
			} else if bucket.waitTicks < bucket.waitTicksUntilPurge {
				bucket.waitTicks++
			}
			bucket.mu.Unlock()
		}
	}(b)
	return b
}

// Inactive Method to get "done" channel
func (b *Bucket) Inactive(ctx context.Context) <-chan struct{} {
	return b.done
}

// CheckLimit Method to check bucket limit
func (b *Bucket) CheckLimit(ctx context.Context) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.requests < b.rateLimit {
		b.requests++
		return true
	}
	return false
}

// ResetLimit Method to reset bucket limit
func (b *Bucket) ResetLimit(ctx context.Context) {
	b.mu.Lock()
	b.requests = 0
	b.waitTicks = 0
	b.mu.Unlock()
}
