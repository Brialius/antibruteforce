package main

import (
	"context"
	"fmt"
	"github.com/Brialius/antibruteforce/internal/config"
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"github.com/Brialius/antibruteforce/internal/grpc"
	"github.com/Brialius/antibruteforce/internal/monitoring"
	"github.com/Brialius/antibruteforce/internal/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func selectStorage(storageType, dsn string) (interfaces.ConfigStorage, error) {
	if storageType == "redis" {
		eventStorage, err := storage.NewRedisStorage(dsn)
		return eventStorage, err
	}
	return nil, errors.Errorf("storage `%s` is not implemented", storageType)
}

var RootCmd = &cobra.Command{
	Use:   "server",
	Short: "Run gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		serverConfig := config.GetGrpcServerConfig()
		storageConfig := config.GetStorageConfig()
		isAbsentParam := false
		if serverConfig.Host == "" {
			isAbsentParam = true
			log.Println("Host is not set")
		}
		if serverConfig.Port == "" {
			isAbsentParam = true
			log.Println("Port is not set")
		}
		if storageConfig.Dsn == "" {
			isAbsentParam = true
			log.Println("Dsn is not set")
		}
		if storageConfig.StorageType == "" {
			isAbsentParam = true
			log.Println("StorageType is not set")
		}
		if isAbsentParam {
			log.Fatal("Some parameters is not set")
		}

		storage, err := selectStorage(storageConfig.StorageType, storageConfig.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer storage.Close(context.Background())

		server := grpc.NewAntiBruteForceServer()
		addr := fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)
		m := &monitoring.PrometheusService{
			Port: serverConfig.MetricsPort,
		}
		log.Printf("Starting monitoring server on %s...", m.Port)
		m.Serve()
		log.Printf("Starting server on %s...", addr)
		err = server.Serve(addr)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cobra.OnInitialize(config.SetConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Config file location")
	RootCmd.PersistentFlags().StringP("metrics-port", "m", "9001", "Port for metrics server")
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("metrics-port", RootCmd.PersistentFlags().Lookup("metrics-port"))
	RootCmd.Flags().StringP("host", "n", "", "host name")
	RootCmd.Flags().IntP("port", "p", 0, "port to listen")
	RootCmd.Flags().StringP("dsn", "d", "", "database connection string")
	RootCmd.Flags().StringP("storage", "s", "", "storage type")
	_ = viper.BindPFlag("grpc-srv-host", RootCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-srv-port", RootCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("dsn", RootCmd.Flags().Lookup("dsn"))
	_ = viper.BindPFlag("storage", RootCmd.Flags().Lookup("storage"))
}

var (
	version = "dev"
	build   = "local"
)

func main() {
	log.Printf("Started antibruteforce gRPC server %s-%s", version, build)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
