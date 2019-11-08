package config

import (
	"github.com/spf13/viper"
	"log"
)

type GrpcClientConfig struct {
	Port     string
	Host     string
	Login    string
	Password string
	Ip       string
}

func GetGrpcClientConfig() *GrpcClientConfig {
	log.Println("Configuring client...")
	viper.SetDefault("grpc-cli-host", "localhost")
	viper.SetDefault("grpc-cli-port", "8080")
	return newGrpcClientConfig()
}

func newGrpcClientConfig() *GrpcClientConfig {
	return &GrpcClientConfig{
		Port:     viper.GetString("grpc-cli-port"),
		Host:     viper.GetString("grpc-cli-host"),
		Login:    viper.GetString("login"),
		Password: viper.GetString("password"),
		Ip:       viper.GetString("ip"),
	}
}
