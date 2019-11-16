package leakybucket

import (
	"context"
	"log"
	"sync"
	"time"
)

type Bucket struct {
	Id                  string
	rateLimit           uint64
	requests            uint64
	waitTicks           uint64
	waitTicksUntilPurge uint64
	done                chan struct{}
	mu                  sync.RWMutex
	once                sync.Once
}

// Might be a parameter
const inactiveCycles = 2

func NewBucket(id string, rateLimit uint64) *Bucket {
	return &Bucket{Id: id, rateLimit: rateLimit, waitTicksUntilPurge: inactiveCycles * rateLimit,
		done: make(chan struct{}, 1), mu: sync.RWMutex{}}
}

func (b *Bucket) Inactive(ctx context.Context) <-chan struct{} {
	return b.done
}

func (b *Bucket) CheckLimit(ctx context.Context) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.requests == 0 {
		b.requests += 1
		b.once.Do(func() {
			go func(bucket *Bucket) {
				ticker := time.NewTicker(time.Duration(uint64(time.Minute) / b.rateLimit))
				defer ticker.Stop()
				for {
					<-ticker.C
					b.mu.RLock()
					if b.waitTicks >= b.waitTicksUntilPurge {
						b.done <- struct{}{}
						b.mu.RUnlock()
						log.Printf("Exit from goroutine for Id: %s", b.Id)
						return
					}
					b.mu.RUnlock()
					b.mu.Lock()
					if b.requests > 0 {
						log.Printf("Id: %s, Requests: %d", b.Id, b.requests)
						b.requests -= 1
						b.waitTicks = 0
					} else {
						if b.waitTicks < b.waitTicksUntilPurge {
							log.Printf("Id: %s, waitTicks: %d", b.Id, b.waitTicks)
							b.waitTicks += 1
						}
					}
					b.mu.Unlock()
				}
			}(b)
		})
		return true
	}
	if b.requests < b.rateLimit {
		b.requests += 1
		return true
	}
	return false
}

func (b *Bucket) ResetLimit(ctx context.Context) {
	b.mu.Lock()
	b.requests = 0
	b.waitTicks = 0
	b.mu.Unlock()
}
