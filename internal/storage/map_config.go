package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"net"
)

type MapConfigStorage struct {
	Whitelist *models.NetList
	Blacklist *models.NetList
}

func NewMapConfigStorage() *MapConfigStorage {
	return &MapConfigStorage{
		Whitelist: &models.NetList{Networks: map[string]*net.IPNet{}},
		Blacklist: &models.NetList{Networks: map[string]*net.IPNet{}},
	}
}

func (m MapConfigStorage) CheckIP(ctx context.Context, ip net.IP) bool {
	for _, value := range m.Whitelist.Networks {
		if value.Contains(ip) {
			return true
		}
	}
	for _, value := range m.Blacklist.Networks {
		if value.Contains(ip) {
			return false
		}
	}
	return true
}

func (m MapConfigStorage) AddToBlackList(ctx context.Context, n *net.IPNet) error {
	m.Blacklist.Networks[n.String()] = n
	return nil
}

func (m MapConfigStorage) DeleteFromBlackList(ctx context.Context, n *net.IPNet) error {
	if _, ok := m.Blacklist.Networks[n.String()]; ok {
		delete(m.Blacklist.Networks, n.String())
		return nil
	}
	return errors.ErrNotFound
}

func (m MapConfigStorage) AddToWhiteList(ctx context.Context, n *net.IPNet) error {
	m.Whitelist.Networks[n.String()] = n
	return nil
}

func (m MapConfigStorage) DeleteFromWhiteList(ctx context.Context, n *net.IPNet) error {
	if _, ok := m.Whitelist.Networks[n.String()]; ok {
		delete(m.Whitelist.Networks, n.String())
		return nil
	}
	return errors.ErrNotFound
}

func (m MapConfigStorage) Close(ctx context.Context) {}
