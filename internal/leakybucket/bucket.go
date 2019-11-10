package leakybucket

import (
	"context"
	"log"
	"sync"
	"time"
)

const waitTicksUntilDone = 10

type Bucket struct {
	Id        string
	rateLimit uint64
	requests  uint64
	waitTicks uint64
	done      chan struct{}
	mu        sync.RWMutex
}

func (b *Bucket) Inactive(ctx context.Context) <-chan struct{} {
	return b.done
}

func NewBucket(id string, rateLimit uint64) *Bucket {
	return &Bucket{Id: id, rateLimit: rateLimit, done: make(chan struct{}, 1), mu: sync.RWMutex{}}
}

func (b *Bucket) CheckLimit(ctx context.Context) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.requests == 0 {
		b.requests += 1
		go func(bucket *Bucket) {
			ticker := time.NewTicker(time.Duration(uint64(time.Minute) / b.rateLimit))
			defer ticker.Stop()
			for {
				<-ticker.C
				b.mu.RLock()
				log.Printf("Id: %s, Requests: %d", b.Id, b.requests)
				if b.waitTicks >= waitTicksUntilDone {
					b.mu.RUnlock()
					b.done <- struct{}{}
					break
				}
				b.mu.RUnlock()
				b.mu.Lock()
				if b.requests > 0 {
					b.requests -= 1
					b.waitTicks = 0
				} else {
					if b.waitTicks < waitTicksUntilDone {
						log.Printf("Id: %s, waitTicks: %d", b.Id, b.waitTicks)
						b.waitTicks += 1
					}
				}
				b.mu.Unlock()
			}
		}(b)
		return true
	}
	if b.requests < b.rateLimit {
		b.requests += 1
		return true
	}
	return false
}

func (b *Bucket) ResetLimit(ctx context.Context, rate uint64) {
	b.mu.Lock()
	b.requests = 0
	b.rateLimit = rate
	b.mu.Unlock()
}
