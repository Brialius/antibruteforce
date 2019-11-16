package main

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	"log"
	"strings"
)

func addToWhiteList(ctx context.Context) {
	if grpcConfig.Ip == "" {
		log.Println("Ip is not set")
	}

	if !strings.Contains(grpcConfig.Ip, "/") {
		log.Printf("`%s` doesn't contain network specificator, seting it for single IP..", grpcConfig.Ip)
		grpcConfig.Ip += netMask
	}

	resp, err := grpcClient.AddToWhiteList(ctx, &api.AddToWhiteListRequest{
		Net: grpcConfig.Ip,
	})

	if err != nil {
		log.Fatal(err)
	}

	if resp.GetError() != "" {
		log.Printf("Error: %v", resp.GetError())
	}
}
