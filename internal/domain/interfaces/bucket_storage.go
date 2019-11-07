package interfaces

import "context"

type BucketStorage interface {
	Close(ctx context.Context)
}
