package services

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"log"
	"net"
)

type AntiBruteForceService struct {
	BucketStorage interfaces.BucketStorage
	ConfigStorage interfaces.ConfigStorage
	LoginLimit    int
	PasswordLimit int
	IpLimit       int
}

func NewAntiBruteForceService(bucketStorage interfaces.BucketStorage, configStorage interfaces.ConfigStorage,
	loginLimit int, passwordLimit int, ipLimit int) *AntiBruteForceService {
	return &AntiBruteForceService{
		BucketStorage: bucketStorage,
		ConfigStorage: configStorage,
		LoginLimit:    loginLimit,
		PasswordLimit: passwordLimit,
		IpLimit:       ipLimit,
	}
}

func (a *AntiBruteForceService) CheckAuth(ctx context.Context, auth *models.Auth) (bool, error) {
	res := a.ConfigStorage.CheckIP(ctx, auth.IpAddr)
	return res, nil
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
