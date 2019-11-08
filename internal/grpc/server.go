package grpc

import "github.com/Brialius/antibruteforce/internal/domain/services"

type AntiBruteForceServer struct {
	AntiBruteForceService *services.AntiBruteForceService
}

func NewAntiBruteForceServer(antiBruteForceService *services.AntiBruteForceService) *AntiBruteForceServer {
	return &AntiBruteForceServer{AntiBruteForceService: antiBruteForceService}
}

func (a *AntiBruteForceServer) Serve(addr string) error {
	return nil
}
