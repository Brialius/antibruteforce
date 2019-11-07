package config

import (
	"github.com/spf13/viper"
	"log"
)

type GrpcServerConfig struct {
	Host        string
	Port        string
	MetricsPort string
}

func GetGrpcServerConfig() *GrpcServerConfig {
	log.Println("Configuring server...")
	viper.SetDefault("grpc-srv-host", "localhost")
	viper.SetDefault("grpc-srv-port", "8080")
	return newGrpcServerConfig()
}

func newGrpcServerConfig() *GrpcServerConfig {
	return &GrpcServerConfig{
		Host:        viper.GetString("grpc-srv-host"),
		Port:        viper.GetString("grpc-srv-port"),
		MetricsPort: viper.GetString("metrics-port"),
	}
}
