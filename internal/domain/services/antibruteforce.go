package services

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"github.com/Brialius/antibruteforce/internal/leakybucket"
	"log"
	"net"
	"time"
)

const (
	duration       = time.Minute
	inactiveCycles = 2
)

type AntiBruteForceService struct {
	BucketStorage interfaces.BucketStorage
	ConfigStorage interfaces.ConfigStorage
	LoginLimit    uint64
	PasswordLimit uint64
	IpLimit       uint64
}

func NewAntiBruteForceService(
	bucketStorage interfaces.BucketStorage, configStorage interfaces.ConfigStorage, loginLimit uint64,
	passwordLimit uint64, ipLimit uint64) *AntiBruteForceService {
	return &AntiBruteForceService{
		BucketStorage: bucketStorage,
		ConfigStorage: configStorage,
		LoginLimit:    loginLimit,
		PasswordLimit: passwordLimit,
		IpLimit:       ipLimit,
	}
}

func (a *AntiBruteForceService) CheckAuth(ctx context.Context, auth *models.Auth) (bool, error) {
	if !a.ConfigStorage.CheckIP(ctx, auth.IpAddr) {
		log.Printf("IP address `%s` is blocked", auth.IpAddr)
		return false, nil
	}
	res := true
	ok, err := a.CheckBucketLimit(ctx, "ip_"+auth.IpAddr.String(), a.IpLimit)
	if err != nil {
		return false, err
	}
	if !ok {
		log.Printf("IP address `%s` requests rate limit is exceeded", auth.IpAddr)
		res = false
	}
	ok, err = a.CheckBucketLimit(ctx, "login_"+auth.Login, a.LoginLimit)
	if err != nil {
		return false, err
	}
	if !ok {
		log.Printf("Login `%s` requests rate limit is exceeded", auth.Login)
		res = false
	}
	ok, err = a.CheckBucketLimit(ctx, "password_"+auth.Password, a.PasswordLimit)
	if err != nil {
		return false, err
	}
	if !ok {
		log.Printf("Password `%s` requests rate limit is exceeded", auth.Password)
		res = false
	}
	return res, nil
}

func (a *AntiBruteForceService) CheckBucketLimit(ctx context.Context, id string, rate uint64) (bool, error) {
	b, err := a.BucketStorage.GetBucket(ctx, id)
	if err != nil {
		if err != errors.ErrBucketNotFound {
			log.Printf("Can't get bucket `%s`: %s", id, err)
			return false, err
		}
		log.Printf("Bucket `%s` doesn't exist, creating..", id)
		if err := a.BucketStorage.SaveBucket(ctx, id, rate,
			leakybucket.NewBucket(id, rate, duration, inactiveCycles)); err != nil {
			log.Printf("Can't create bucket `%s`: %s", id, err)
			return false, err
		}
		b, err = a.BucketStorage.GetBucket(ctx, id)
		if err != nil {
			return false, err
		}
	}
	return b.CheckLimit(ctx), err
}

func (a *AntiBruteForceService) AddToWhiteList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Adding %s to whitelist", n)
	err := a.ConfigStorage.AddToWhiteList(ctx, n)
	if err != nil {
		log.Printf("Adding %s to whitelist is failed: %s", n, err)
	}
	return err
}

func (a *AntiBruteForceService) AddToBlackList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Adding %s to blacklist", n)
	err := a.ConfigStorage.AddToBlackList(ctx, n)
	if err != nil {
		log.Printf("Adding %s to blacklist is failed: %s", n, err)
	}
	return err
}

func (a *AntiBruteForceService) DeleteFromWhiteList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Deleting %s from whitelist", n)
	err := a.ConfigStorage.DeleteFromWhiteList(ctx, n)
	if err != nil {
		log.Printf("Deleting %s from whitelist is failed: %s", n, err)
	}
	return err
}

func (a *AntiBruteForceService) DeleteFromBlackList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Deleting %s from blacklist", n)
	err := a.ConfigStorage.DeleteFromBlackList(ctx, n)
	if err != nil {
		log.Printf("Deleting %s from blacklist is failed: %s", n, err)
	}
	return err
}

func (a *AntiBruteForceService) ResetLimit(ctx context.Context, login string, n *net.IPNet) error {
	panic("implement me")
}
