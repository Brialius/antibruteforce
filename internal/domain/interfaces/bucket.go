package interfaces

import "context"

type Bucket interface {
	CheckLimit(ctx context.Context) bool
	ResetLimit(ctx context.Context)
	Inactive(ctx context.Context) <-chan struct{}
}
