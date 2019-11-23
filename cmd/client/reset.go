package main

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	"log"
)

func resetLimit(ctx context.Context) {
	isAbsentParam := false

	if grpcConfig.Login == "" {
		isAbsentParam = true
		log.Println("Login is not set")
	}

	if grpcConfig.Ip == "" {
		isAbsentParam = true
		log.Println("Ip is not set")
	}

	if isAbsentParam {
		log.Fatal("Some parameters is not set")
	}

	resp, err := grpcClient.ResetLimit(ctx, &api.ResetLimitRequest{
		Login: grpcConfig.Login,
		Ip:    grpcConfig.Ip,
	})

	if err != nil {
		log.Fatal(err)
	}

	if resp.GetError() != "" {
		log.Printf("Reset error: %v", resp.GetError())
	}
}
