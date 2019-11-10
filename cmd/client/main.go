package main

import (
	"context"
	"fmt"
	"github.com/Brialius/antibruteforce/internal/config"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ReqTimeout = time.Second * 10

func newGrpcClient(ctx context.Context, host string, port string) api.AntiBruteForceServiceClient {
	server := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.DialContext(ctx, server, grpc.WithInsecure(), grpc.WithUserAgent("antibruteforce client"))
	if err != nil {
		log.Fatal(err)
	}
	return api.NewAntiBruteForceServiceClient(conn)
}

func checkAuth(ctx context.Context) {
	resp, err := grpcClient.CheckAuth(ctx, &api.CheckAuthRequest{
		Auth: &api.Auth{
			Login:    grpcConfig.Login,
			Password: grpcConfig.Password,
			Ip:       grpcConfig.Ip,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Check error: %v", resp.GetError())
	fmt.Printf("Check result: %v", resp.GetOk())
}

func addToWhiteList(ctx context.Context) {
	panic("implement me")
}

func addToBlackList(ctx context.Context) {
	panic("implement me")
}

func deleteFromWhiteList(ctx context.Context) {
	panic("implement me")
}

func deleteFromBlackList(ctx context.Context) {
	panic("implement me")
}

func resetLimit(ctx context.Context) {
	panic("implement me")
}

var RootCmd = &cobra.Command{
	Use:       "client [check, reset, allow, disallow, block, unblock]",
	Short:     "Run gRPC client",
	ValidArgs: []string{"check", "reset", "allow", "disallow", "block", "unblock"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		grpcConfig = config.GetGrpcClientConfig()
		ctx, cancel := context.WithTimeout(context.Background(), ReqTimeout)
		grpcClient = newGrpcClient(ctx, grpcConfig.Host, grpcConfig.Port)
		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGINT)
			<-stop
			log.Printf("Interrupt signal")
			cancel()
		}()
		switch args[0] {
		case "check":
			checkAuth(ctx)
		case "reset":
			resetLimit(ctx)
		case "allow":
			addToWhiteList(ctx)
		case "disallow":
			deleteFromWhiteList(ctx)
		case "block":
			addToBlackList(ctx)
		case "unblock":
			deleteFromBlackList(ctx)
		}
	},
}

func init() {
	cobra.OnInitialize(config.SetConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Config file location")
	RootCmd.Flags().StringP("login", "l", "", "login")
	RootCmd.Flags().StringP("password", "w", "", "password")
	RootCmd.Flags().StringP("ip", "i", "", "ip address")
	RootCmd.Flags().StringP("host", "n", "", "host name")
	RootCmd.Flags().IntP("port", "p", 0, "port to listen")
	// bind flags to viper
	_ = viper.BindPFlag("grpc-cli-host", RootCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-cli-port", RootCmd.Flags().Lookup("port"))
}

var (
	version    = "dev"
	build      = "local"
	grpcConfig *config.GrpcClientConfig
	grpcClient api.AntiBruteForceServiceClient
)

func main() {
	log.Printf("Started antibruteforce gRPC client %s-%s", version, build)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
