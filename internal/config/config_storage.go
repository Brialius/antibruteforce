package config

import (
	"github.com/spf13/viper"
	"log"
)

// StorageConfig Struct for config storage
type StorageConfig struct {
	Dsn         string
	StorageType string
}

// GetStorageConfig Get storage config
func GetStorageConfig() *StorageConfig {
	log.Println("Configuring storage...")
	viper.SetDefault("dsn", "")
	viper.SetDefault("storage", "map")
	return newDbConfig()
}

func newDbConfig() *StorageConfig {
	return &StorageConfig{
		Dsn:         viper.GetString("dsn"),
		StorageType: viper.GetString("storage"),
	}
}
