package storage

import (
	"context"
	"fmt"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	bolt "github.com/coreos/bbolt"
	"log"
	"net"
	"time"
)

const (
	whitelistBucket = "whitelist"
	blacklistBucket = "blacklist"
)

// BoltConfigStorage bolt DB storage for config
type BoltConfigStorage struct {
	db *bolt.DB
}

func (b *BoltConfigStorage) AddToBlackList(ctx context.Context, n *net.IPNet) error {
	return b.addToBucket(n.String(), "", blacklistBucket)
}

func (b *BoltConfigStorage) DeleteFromBlackList(ctx context.Context, n *net.IPNet) error {
	return b.deleteFromBucket(n.String(), blacklistBucket)
}

func (b *BoltConfigStorage) AddToWhiteList(ctx context.Context, n *net.IPNet) error {
	return b.addToBucket(n.String(), "", whitelistBucket)
}

func (b *BoltConfigStorage) DeleteFromWhiteList(ctx context.Context, n *net.IPNet) error {
	return b.deleteFromBucket(n.String(), whitelistBucket)
}

func (b *BoltConfigStorage) CheckIP(ctx context.Context, ip net.IP) (bool, error) {
	whitelisted, err := b.containsInList(ip, whitelistBucket)
	if err != nil {
		return false, err
	}
	if whitelisted {
		return true, nil
	}

	blacklisted, err := b.containsInList(ip, blacklistBucket)
	if err != nil {
		return false, err
	}
	if blacklisted {
		return false, nil
	}
	return true, nil
}

func (b *BoltConfigStorage) containsInList(ip net.IP, bucket string) (bool, error) {
	err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		return bkt.ForEach(func(k, v []byte) error {
			_, n, _ := net.ParseCIDR(string(k))
			if n.Contains(ip) {
				// Send an error to break the ForEach loop and exit from function
				return fmt.Errorf("net contains ip")
			}
			return nil
		})
	})
	if err != nil {
		return true, nil
	}
	return false, nil
}

func (b *BoltConfigStorage) Close(ctx context.Context) error {
	return b.db.Close()
}

func (b *BoltConfigStorage) addToBucket(key, value, bucket string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		e := bucket.Put([]byte(key), []byte(value))
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		log.Printf("can't add key %s to bolt db: %s", key, err)
	}
	return err
}

func (b *BoltConfigStorage) deleteFromBucket(key, bucket string) error {
	if !b.isExist(key, bucket) {
		log.Printf("key `%s` doesn't exist in bucket `%s`", key, bucket)
		return errors.ErrNotFound
	}
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		e := b.Delete([]byte(key))
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		log.Printf("can't delete key `%s` from bucket `%s`: %s", key, bucket, err)
	}
	return err
}

func (b *BoltConfigStorage) isExist(key, bucket string) bool {
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket.Get([]byte(key)) == nil {
			return errors.ErrNotFound
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

func NewBoltConfigStorage(timeout time.Duration, fileName string) (*BoltConfigStorage, error) {
	db, err := bolt.Open(fileName, 0600, &bolt.Options{Timeout: timeout})
	if err != nil {
		log.Printf("can't open bolt db: %s", err)
		return nil, err
	}

	buckets := []string{whitelistBucket, blacklistBucket}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bktName := range buckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(bktName)); e != nil {
				log.Printf("failed to create top level bucket %s: %s", bktName, err)
				return e
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("bolt store created: %s", fileName)
	return &BoltConfigStorage{db: db}, nil
}
