package models

import "github.com/Brialius/antibruteforce/internal/domain/interfaces"

type Bucket struct {
	Id        string
	RateLimit int //requests per minute
	Storage   *interfaces.BucketStorage
}
