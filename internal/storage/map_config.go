package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"net"
	"sync"
)

type MapConfigStorage struct {
	Whitelist   *models.NetList
	Blacklist   *models.NetList
	whitelistMu sync.RWMutex
	blacklistMu sync.RWMutex
}

func NewMapConfigStorage() *MapConfigStorage {
	return &MapConfigStorage{
		Whitelist:   &models.NetList{Networks: map[string]*net.IPNet{}},
		Blacklist:   &models.NetList{Networks: map[string]*net.IPNet{}},
		whitelistMu: sync.RWMutex{},
		blacklistMu: sync.RWMutex{},
	}
}

func (m MapConfigStorage) CheckIP(ctx context.Context, ip net.IP) bool {
	m.whitelistMu.RLock()
	defer m.whitelistMu.RUnlock()

	for _, value := range m.Whitelist.Networks {
		if value.Contains(ip) {
			return true
		}
	}

	m.blacklistMu.RLock()
	defer m.blacklistMu.RUnlock()

	for _, value := range m.Blacklist.Networks {
		if value.Contains(ip) {
			return false
		}
	}
	return true
}

func (m MapConfigStorage) AddToBlackList(ctx context.Context, n *net.IPNet) error {
	m.blacklistMu.Lock()
	m.Blacklist.Networks[n.String()] = n
	m.blacklistMu.Unlock()
	return nil
}

func (m MapConfigStorage) DeleteFromBlackList(ctx context.Context, n *net.IPNet) error {
	m.blacklistMu.Lock()
	defer m.blacklistMu.Unlock()

	if _, ok := m.Blacklist.Networks[n.String()]; ok {
		delete(m.Blacklist.Networks, n.String())
		return nil
	}
	return errors.ErrNotFound
}

func (m MapConfigStorage) AddToWhiteList(ctx context.Context, n *net.IPNet) error {
	m.whitelistMu.Lock()
	m.Whitelist.Networks[n.String()] = n
	m.whitelistMu.Unlock()
	return nil
}

func (m MapConfigStorage) DeleteFromWhiteList(ctx context.Context, n *net.IPNet) error {
	m.whitelistMu.Lock()
	defer m.whitelistMu.Unlock()

	if _, ok := m.Whitelist.Networks[n.String()]; ok {
		delete(m.Whitelist.Networks, n.String())
		return nil
	}
	return errors.ErrNotFound
}

func (m MapConfigStorage) Close(ctx context.Context) {}
