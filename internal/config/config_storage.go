package config

import (
	"github.com/spf13/viper"
	"log"
)

type StorageConfig struct {
	Dsn         string
	StorageType string
}

func GetStorageConfig() *StorageConfig {
	log.Println("Configuring storage...")
	viper.SetDefault("dsn", "")
	viper.SetDefault("storage", "")
	return newDbConfig()
}

func newDbConfig() *StorageConfig {
	return &StorageConfig{
		Dsn:         viper.GetString("dsn"),
		StorageType: viper.GetString("storage"),
	}
}
