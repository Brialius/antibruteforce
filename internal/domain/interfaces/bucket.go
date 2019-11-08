package interfaces

import "context"

type Bucket interface {
	CheckLimit(ctx context.Context, id string) bool
	ResetLimit(ctx context.Context, id string)
}
