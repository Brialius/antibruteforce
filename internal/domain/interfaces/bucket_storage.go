package interfaces

import (
	"context"
)

type BucketStorage interface {
	SaveBucket(ctx context.Context, id string, rateLimit uint64, bucket Bucket) error
	DeleteBucket(ctx context.Context, id string) error
	GetBucket(ctx context.Context, id string) (Bucket, error)
	Close(ctx context.Context)
}
