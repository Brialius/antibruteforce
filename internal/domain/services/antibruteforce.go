package services

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"github.com/Brialius/antibruteforce/internal/domain/models"
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
	//TODO: Implement
	return false, nil
}
