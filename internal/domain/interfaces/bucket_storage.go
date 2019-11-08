package interfaces

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/models"
)

type BucketStorage interface {
	CreateBucket(ctx context.Context, id string, rateLimit int) (*models.Bucket, error)
	DeleteBucket(ctx context.Context, id string, rateLimit int) error
	GetBucket(ctx context.Context, id string, rateLimit int) (*models.Bucket, error)
	Close(ctx context.Context)
}
