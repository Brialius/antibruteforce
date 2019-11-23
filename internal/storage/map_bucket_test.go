package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapBucketStorage(t *testing.T) {
	m := NewMapBucketStorage()

	ctx := context.Background()

	t.Run("Check empty storage", func(t *testing.T) {
		assert.Equal(t, 0, len(m.Buckets))
	})

	t.Run("Add bucket", func(t *testing.T) {
		err := m.SaveBucket(ctx, "testAdd", 0, bucketMock{})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, 1, len(m.Buckets))
	})

	t.Run("Get bucket", func(t *testing.T) {
		b, err := m.GetBucket(ctx, "testAdd")
		if err != nil {
			t.Error(err)
		}
		assert.NotNil(t, b)
	})

	t.Run("Delete bucket", func(t *testing.T) {
		err := m.DeleteBucket(ctx, "testAdd")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, 0, len(m.Buckets))
	})

	t.Run("Delete non-exist bucket", func(t *testing.T) {
		err := m.DeleteBucket(ctx, "testAdd")
		assert.NotNil(t, err)
	})
}

// Mocks
type bucketMock struct {
	Name string
}

func (b bucketMock) CheckLimit(ctx context.Context) bool {
	return true
}

func (b bucketMock) ResetLimit(ctx context.Context) {}

func (b bucketMock) Inactive(ctx context.Context) <-chan struct{} {
	return nil
}
