package storage

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
	"time"
)

func TestBoltConfigStorage_CheckIP(t *testing.T) {
	ctx := context.Background()
	boltFile := "testBoltDb"
	c, err := NewBoltConfigStorage(10*time.Second, boltFile)
	defer func() {
		if err := os.Remove(boltFile); err != nil {
			t.Fatalf("can't delete bolt db file: %s", err)
		}
	}()
	defer func() {
		if err := c.Close(ctx); err != nil {
			t.Fatalf("can't close bolt db connection: %s", err)
		}
	}()

	if err != nil {
		t.Errorf("can't create bolt storage: %s", err)
	}
	t.Run("Check positive", func(t *testing.T) {
		res, err := c.CheckIP(ctx, net.ParseIP("11.11.11.11"))
		assert.Equal(t, true, res)
		assert.Nil(t, err)
	})

	t.Run("Add to whitelist", func(t *testing.T) {
		err := c.AddToWhiteList(ctx, getNet(t, "10.0.0.0/8"))
		assert.Nil(t, err)
		err = c.AddToWhiteList(ctx, getNet(t, "192.168.11.11/32"))
		assert.True(t, c.isExist("10.0.0.0/8", whitelistBucket))
		assert.True(t, c.isExist("192.168.11.11/32", whitelistBucket))
		assert.Nil(t, err)
	})

	t.Run("Add to blacklist", func(t *testing.T) {
		err := c.AddToBlackList(ctx, getNet(t, "192.168.0.0/16"))
		assert.Nil(t, err)
		err = c.AddToBlackList(ctx, getNet(t, "10.11.11.11/32"))
		assert.Nil(t, err)
		assert.True(t, c.isExist("192.168.0.0/16", blacklistBucket))
		assert.True(t, c.isExist("10.11.11.11/32", blacklistBucket))
	})

	t.Run("Check whitelisted", func(t *testing.T) {
		res, err := c.CheckIP(ctx, net.ParseIP("10.10.10.10"))
		assert.Equal(t, true, res)
		assert.Nil(t, err)
	})

	t.Run("Check not listed", func(t *testing.T) {
		res, err := c.CheckIP(ctx, net.ParseIP("11.11.11.11"))
		assert.Nil(t, err)
		assert.Equal(t, true, res)
	})

	t.Run("Check whitelisted and blacklisted", func(t *testing.T) {
		res, err := c.CheckIP(ctx, net.ParseIP("10.11.11.11"))
		assert.Nil(t, err)
		assert.Equal(t, true, res)
	})

	t.Run("Check blacklisted and whitelisted", func(t *testing.T) {
		res, err := c.CheckIP(ctx, net.ParseIP("192.168.11.11"))
		assert.Nil(t, err)
		assert.Equal(t, true, res)
	})

	t.Run("Check blacklisted", func(t *testing.T) {
		res, err := c.CheckIP(ctx, net.ParseIP("192.168.10.10"))
		assert.Nil(t, err)
		assert.Equal(t, false, res)
	})

	t.Run("Delete from whitelist", func(t *testing.T) {
		err := c.DeleteFromWhiteList(ctx, getNet(t, "10.0.0.0/8"))
		assert.False(t, c.isExist("10.0.0.0/8", whitelistBucket))
		assert.NoError(t, err)
		err = c.DeleteFromWhiteList(ctx, getNet(t, "192.168.11.11/32"))
		assert.False(t, c.isExist("192.168.11.11/32", whitelistBucket))
		assert.NoError(t, err)
	})

	t.Run("Delete from blacklist", func(t *testing.T) {
		err := c.DeleteFromBlackList(ctx, getNet(t, "192.168.0.0/16"))
		assert.NoError(t, err)
		err = c.DeleteFromBlackList(ctx, getNet(t, "10.11.11.11/32"))
		assert.NoError(t, err)
		assert.False(t, c.isExist("192.168.0.0/16", blacklistBucket))
		assert.False(t, c.isExist("10.11.11.11/32", blacklistBucket))
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
