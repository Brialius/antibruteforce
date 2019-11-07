package storage

import "context"

type RedisStorage struct {
}

func NewRedisStorage(dsn string) (*RedisStorage, error) {
	return &RedisStorage{}, nil
}

func (r *RedisStorage) Close(ctx context.Context) {

}
