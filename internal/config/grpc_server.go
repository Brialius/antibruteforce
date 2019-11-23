package config

import (
	"github.com/spf13/viper"
	"log"
)

// GrpcServerConfig gRPC server configuration struct
type GrpcServerConfig struct {
	Host          string
	Port          string
	MetricsPort   string
	LoginLimit    uint64
	PasswordLimit uint64
	IPLimit       uint64
}

// GetGrpcServerConfig Get gRPC server configuration
func GetGrpcServerConfig() *GrpcServerConfig {
	log.Println("Configuring server...")
	viper.SetDefault("grpc-srv-host", "localhost")
	viper.SetDefault("grpc-srv-port", "8080")
	viper.SetDefault("login-limit", 10)
	viper.SetDefault("password-limit", 100)
	viper.SetDefault("ip-limit", 1000)
	return newGrpcServerConfig()
}

func newGrpcServerConfig() *GrpcServerConfig {
	return &GrpcServerConfig{
		Host:          viper.GetString("grpc-srv-host"),
		Port:          viper.GetString("grpc-srv-port"),
		MetricsPort:   viper.GetString("metrics-port"),
		LoginLimit:    viper.GetUint64("login-limit"),
		PasswordLimit: viper.GetUint64("password-limit"),
		IPLimit:       viper.GetUint64("ip-limit"),
	}
}
