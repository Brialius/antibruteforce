package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"log"
	"sync"
)

// MapBucketStorage Map bucket storage struct
type MapBucketStorage struct {
	Buckets map[string]interfaces.Bucket
	mu      *sync.RWMutex
}

// SaveBucket Method to save bucket into storage
func (m MapBucketStorage) SaveBucket(ctx context.Context, id string, rateLimit uint64, bucket interfaces.Bucket) error {
	m.mu.Lock()
	m.Buckets[id] = bucket
	m.mu.Unlock()

	go func(ctx context.Context, bucketStorage MapBucketStorage, id string, done <-chan struct{}) {
		<-done
		err := bucketStorage.DeleteBucket(ctx, id)
		if err != nil {
			log.Printf("Can't delete inactive bucket %s: %s", id, err)
		}
		log.Printf("Deleted inactive bucket %s", id)
	}(ctx, m, id, bucket.Inactive(ctx))
	return nil
}

// DeleteBucket Method to delete bucket from storage
func (m MapBucketStorage) DeleteBucket(ctx context.Context, id string) error {
	_, err := m.GetBucket(ctx, id)
	if err != nil {
		return err
	}

	m.mu.Lock()
	delete(m.Buckets, id)
	m.mu.Unlock()
	return nil
}

// GetBucket Method to get bucket from storage
func (m MapBucketStorage) GetBucket(ctx context.Context, id string) (interfaces.Bucket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if res, ok := m.Buckets[id]; ok {
		return res, nil
	}
	return nil, errors.ErrBucketNotFound
}

// Close Method to close connection to storage (to implement Bucket storage interface)
func (m MapBucketStorage) Close(ctx context.Context) {}

// NewMapBucketStorage Map bucket storage constructor
func NewMapBucketStorage() *MapBucketStorage {
	return &MapBucketStorage{
		Buckets: map[string]interfaces.Bucket{},
		mu:      &sync.RWMutex{},
	}
}
