package main

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	"log"
	"strings"
)

func deleteFromWhiteList(ctx context.Context) {
	if grpcConfig.IP == "" {
		log.Println("IP is not set")
	}

	if !strings.Contains(grpcConfig.IP, "/") {
		log.Printf("`%s` doesn't contain network specificator, seting it for single IP..", grpcConfig.IP)
		grpcConfig.IP += netMask
	}

	resp, err := grpcClient.DeleteFromWhiteList(ctx, &api.DeleteFromWhiteListRequest{
		Net: grpcConfig.IP,
	})

	if err != nil {
		log.Fatal(err)
	}

	if resp.GetError() != "" {
		log.Printf("Error: %v", resp.GetError())
	}
}
