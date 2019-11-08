package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"sync"
)

type MapBucketStorage struct {
	Buckets map[string]*interfaces.Bucket
	Pool    sync.Pool
}

func (m MapBucketStorage) CreateBucket(ctx context.Context, id string, rateLimit int) (*models.Bucket, error) {
	panic("implement me")
}

func (m MapBucketStorage) DeleteBucket(ctx context.Context, id string) error {
	panic("implement me")
}

func (m MapBucketStorage) GetBucket(ctx context.Context, id string) (*models.Bucket, error) {
	panic("implement me")
}

func (m MapBucketStorage) Close(ctx context.Context) {
	panic("implement me")
}

func NewMapBucketStorage() *MapBucketStorage {
	return &MapBucketStorage{
		Buckets: map[string]*interfaces.Bucket{},
	}
}
