package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMapConfigStorage_CheckIP(t *testing.T) {
	ctx := context.Background()
	c := NewMapConfigStorage()

	t.Run("Check positive", func(t *testing.T) {
		res := c.CheckIP(ctx, net.ParseIP("11.11.11.11"))
		assert.Equal(t, true, res)
	})

	t.Run("Add to whitelist", func(t *testing.T) {
		_ = c.AddToWhiteList(ctx, getNet(t, "10.0.0.0/8"))
		_ = c.AddToWhiteList(ctx, getNet(t, "192.168.11.11/32"))
		assert.Contains(t, c.Whitelist.Networks, "10.0.0.0/8")
		assert.Contains(t, c.Whitelist.Networks, "192.168.11.11/32")
	})

	t.Run("Add to blacklist", func(t *testing.T) {
		_ = c.AddToBlackList(ctx, getNet(t, "192.168.0.0/16"))
		_ = c.AddToBlackList(ctx, getNet(t, "10.11.11.11/32"))
		assert.Contains(t, c.Blacklist.Networks, "192.168.0.0/16")
		assert.Contains(t, c.Blacklist.Networks, "10.11.11.11/32")
	})

	t.Run("Check whitelisted", func(t *testing.T) {
		res := c.CheckIP(ctx, net.ParseIP("10.10.10.10"))
		assert.Equal(t, true, res)
	})

	t.Run("Check not listed", func(t *testing.T) {
		res := c.CheckIP(ctx, net.ParseIP("11.11.11.11"))
		assert.Equal(t, true, res)
	})

	t.Run("Check whitelisted and blacklisted", func(t *testing.T) {
		res := c.CheckIP(ctx, net.ParseIP("10.11.11.11"))
		assert.Equal(t, true, res)
	})

	t.Run("Check blacklisted and whitelisted", func(t *testing.T) {
		res := c.CheckIP(ctx, net.ParseIP("192.168.11.11"))
		assert.Equal(t, true, res)
	})

	t.Run("Check blacklisted", func(t *testing.T) {
		res := c.CheckIP(ctx, net.ParseIP("192.168.10.10"))
		assert.Equal(t, false, res)
	})

	t.Run("Delete from whitelist", func(t *testing.T) {
		err := c.DeleteFromWhiteList(ctx, getNet(t, "10.0.0.0/8"))
		assert.NotContains(t, c.Whitelist.Networks, "10.0.0.0/8")
		assert.NoError(t, err)
		err = c.DeleteFromWhiteList(ctx, getNet(t, "192.168.11.11/32"))
		assert.NotContains(t, c.Whitelist.Networks, "192.168.11.11/32")
		assert.NoError(t, err)
	})

	t.Run("Delete from blacklist", func(t *testing.T) {
		err := c.DeleteFromBlackList(ctx, getNet(t, "192.168.0.0/16"))
		assert.NoError(t, err)
		err = c.DeleteFromBlackList(ctx, getNet(t, "10.11.11.11/32"))
		assert.NoError(t, err)
		assert.NotContains(t, c.Blacklist.Networks, "192.168.0.0/16")
		assert.NotContains(t, c.Blacklist.Networks, "10.11.11.11/32")
	})

	t.Run("Delete from whitelist negative", func(t *testing.T) {
		err := c.DeleteFromWhiteList(ctx, getNet(t, "10.0.0.0/8"))
		assert.EqualError(t, err, errors.ErrNotFound.Error())
		err = c.DeleteFromWhiteList(ctx, getNet(t, "192.168.11.11/32"))
		assert.EqualError(t, err, errors.ErrNotFound.Error())
	})

	t.Run("Delete from blacklist negative", func(t *testing.T) {
		err := c.DeleteFromBlackList(ctx, getNet(t, "192.168.0.0/16"))
		assert.EqualError(t, err, errors.ErrNotFound.Error())
		err = c.DeleteFromBlackList(ctx, getNet(t, "10.11.11.11/32"))
		assert.EqualError(t, err, errors.ErrNotFound.Error())
	})

}

func getNet(t *testing.T, s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		t.Errorf("can't parse provided network: %s", err)
	}
	return n
}
