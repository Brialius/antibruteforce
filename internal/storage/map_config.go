package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"net"
)

type MapConfigStorage struct {
	Whitelist *models.NetList
	Blacklist *models.NetList
}

func NewMapConfigStorage() *MapConfigStorage {
	return &MapConfigStorage{
		Whitelist: &models.NetList{},
		Blacklist: &models.NetList{},
	}
}

func (m MapConfigStorage) AddToBlackList(ctx context.Context, net net.IPNet) {
	panic("implement me")
}

func (m MapConfigStorage) DeleteFromBlackList(ctx context.Context, net net.IPNet) {
	panic("implement me")
}

func (m MapConfigStorage) AddToWhiteList(ctx context.Context, net net.IPNet) {
	panic("implement me")
}

func (m MapConfigStorage) DeleteFromWhiteList(ctx context.Context, net net.IPNet) {
	panic("implement me")
}

func (m MapConfigStorage) Close(ctx context.Context) {
	panic("implement me")
}
