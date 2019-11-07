package interfaces

import "context"

type ConfigStorage interface {
	Close(ctx context.Context)
}
