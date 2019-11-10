package leakybucket

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"log"
)

type LeakyBucket struct {
	bucketStorage interfaces.BucketStorage
}

func NewLeakyBucket(bucketStorage interfaces.BucketStorage) *LeakyBucket {
	return &LeakyBucket{bucketStorage: bucketStorage}
}

func (l *LeakyBucket) CheckBucketLimit(ctx context.Context, id string, rate uint64) (bool, error) {
	b, err := l.bucketStorage.GetBucket(ctx, id)
	if err != nil {
		if err != errors.ErrBucketNotFound {
			log.Printf("Can't get bucket `%s`: %s", id, err)
			return false, err
		}
		log.Printf("Bucket `%s` doesn't exist, creating..", id)
		if err := l.bucketStorage.CreateBucket(ctx, id, rate, NewBucket(id, rate)); err != nil {
			log.Printf("Can't create bucket `%s`: %s", id, err)
			return false, err
		}
		b, err = l.bucketStorage.GetBucket(ctx, id)
		if err != nil {
			return false, err
		}
	}
	return b.CheckLimit(ctx), err
}
