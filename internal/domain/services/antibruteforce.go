package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"github.com/Brialius/antibruteforce/internal/config"
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

// AntiBruteForceService struct
type AntiBruteForceService struct {
	BucketStorage interfaces.BucketStorage
	ConfigStorage interfaces.ConfigStorage
	LoginLimit    uint64
	PasswordLimit uint64
	IPLimit       uint64
}

// NewAntiBruteForceService Constructor for AntiBruteForceService
func NewAntiBruteForceService(
	bucketStorage interfaces.BucketStorage, configStorage interfaces.ConfigStorage, loginLimit uint64,
	passwordLimit uint64, ipLimit uint64) *AntiBruteForceService {
	return &AntiBruteForceService{
		BucketStorage: bucketStorage,
		ConfigStorage: configStorage,
		LoginLimit:    loginLimit,
		PasswordLimit: passwordLimit,
		IPLimit:       ipLimit,
	}
}

// CheckAuth Method to check limits for Auth
func (a *AntiBruteForceService) CheckAuth(ctx context.Context, auth *models.Auth) (bool, error) {
	ok, err := a.ConfigStorage.CheckIP(ctx, auth.IPAddr)
	if err != nil {
		log.Printf("Can't check IP address `%s`: %s", auth.IPAddr, err)
		return false, err
	}

	if !ok {
		if config.Verbose {
			log.Printf("IP address `%s` is blocked", auth.IPAddr)
		}
		return false, nil
	}
	res := true

	ok, err = a.CheckBucketLimit(ctx, "ip_"+auth.IPAddr.String(), a.IPLimit)
	if err != nil {
		return false, err
	}
	if !ok {
		if config.Verbose {
			log.Printf("IP address `%s` requests rate limit is exceeded", auth.IPAddr)
		}
		res = false
	}
	ok, err = a.CheckBucketLimit(ctx, "login_"+auth.Login, a.LoginLimit)
	if err != nil {
		return false, err
	}
	if !ok {
		if config.Verbose {
			log.Printf("Login `%s` requests rate limit is exceeded", auth.Login)
		}
		res = false
	}

	h := sha256.New()
	if _, err := h.Write([]byte(auth.Password)); err != nil {
		log.Printf("Hash counting error: %s", err)
		return false, err
	}
	hashedPwd := base64.URLEncoding.EncodeToString(h.Sum([]byte(nil)))

	ok, err = a.CheckBucketLimit(ctx, "password_"+hashedPwd, a.PasswordLimit)
	if err != nil {
		return false, err
	}
	if !ok {
		log.Printf("Password `%s` requests rate limit is exceeded", hashedPwd)
		res = false
	}
	return res, nil
}

// CheckBucketLimit Method to check limit for particular bucket
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

// AddToWhiteList Method to add IP address to Whitelist
func (a *AntiBruteForceService) AddToWhiteList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Adding %s to whitelist", n)
	err := a.ConfigStorage.AddToWhiteList(ctx, n)
	if err != nil {
		log.Printf("Adding %s to whitelist is failed: %s", n, err)
	}
	return err
}

// AddToBlackList Method to add IP address to Blacklist
func (a *AntiBruteForceService) AddToBlackList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Adding %s to blacklist", n)
	err := a.ConfigStorage.AddToBlackList(ctx, n)
	if err != nil {
		log.Printf("Adding %s to blacklist is failed: %s", n, err)
	}
	return err
}

// DeleteFromWhiteList Method to delete IP address from Whitelist
func (a *AntiBruteForceService) DeleteFromWhiteList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Deleting %s from whitelist", n)
	err := a.ConfigStorage.DeleteFromWhiteList(ctx, n)
	if err != nil {
		log.Printf("Deleting %s from whitelist is failed: %s", n, err)
	}
	return err
}

// DeleteFromBlackList Method to add IP address from Blacklist
func (a *AntiBruteForceService) DeleteFromBlackList(ctx context.Context, n *net.IPNet) error {
	log.Printf("Deleting %s from blacklist", n)
	err := a.ConfigStorage.DeleteFromBlackList(ctx, n)
	if err != nil {
		log.Printf("Deleting %s from blacklist is failed: %s", n, err)
	}
	return err
}

// ResetLimit Method to reset limits for login/IP address pair
func (a *AntiBruteForceService) ResetLimit(ctx context.Context, login string, ip *net.IP) error {
	b, err := a.BucketStorage.GetBucket(ctx, "login_"+login)
	if err != nil {
		if err != errors.ErrBucketNotFound {
			log.Printf("Can't get bucket `%s`: %s", "login_"+login, err)
			return err
		}
		log.Printf("Bucket `%s` doesn't exist", "login_"+login)
		return err
	}
	b.ResetLimit(ctx)
	log.Printf("Limit for login `%s` has been reset ", login)

	b, err = a.BucketStorage.GetBucket(ctx, "ip_"+ip.String())
	if err != nil {
		if err != errors.ErrBucketNotFound {
			log.Printf("Can't get bucket `%s`: %s", "ip_"+ip.String(), err)
			return err
		}
		log.Printf("Bucket `%s` doesn't exist", "ip_"+ip.String())
		return err
	}
	b.ResetLimit(ctx)
	log.Printf("Limit for IP `%s` has been reset ", ip.String())
	return nil
}
