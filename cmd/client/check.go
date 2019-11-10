package main

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	"log"
)

func checkAuth(ctx context.Context) {
	isAbsentParam := false

	if grpcConfig.Login == "" {
		isAbsentParam = true
		log.Println("Login is not set")
	}

	if grpcConfig.Password == "" {
		isAbsentParam = true
		log.Println("Password is not set")
	}

	if grpcConfig.Ip == "" {
		isAbsentParam = true
		log.Println("Ip is not set")
	}

	if isAbsentParam {
		log.Fatal("Some parameters is not set")
	}

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

	if resp.GetError() != "" {
		log.Printf("Check error: %v", resp.GetError())
	}
	log.Printf("Check result: %v", resp.GetOk())
}
