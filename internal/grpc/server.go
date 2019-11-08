package grpc

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"github.com/Brialius/antibruteforce/internal/domain/services"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type AntiBruteForceServer struct {
	AntiBruteForceService *services.AntiBruteForceService
}

func (a *AntiBruteForceServer) CheckAuth(ctx context.Context, req *api.CheckAuthRequest) (*api.CheckAuthResponse, error) {
	log.Printf("Hello there")
	ip := net.ParseIP(req.GetAuth().GetIp())
	ok, err := a.AntiBruteForceService.CheckAuth(ctx, &models.Auth{
		Login:    req.GetAuth().GetLogin(),
		Password: req.GetAuth().GetPassword(),
		IpAddr:   ip,
	})
	if err != nil {
		return &api.CheckAuthResponse{
			Result: &api.CheckAuthResponse_Error{Error: "check error"},
		}, nil
	}
	return &api.CheckAuthResponse{
		Result: &api.CheckAuthResponse_Ok{Ok: ok},
	}, nil
}

func (a *AntiBruteForceServer) AddToWhiteList(context.Context, *api.AddToWhiteListRequest) (*api.AddToWhiteListResponse, error) {
	panic("implement me")
}

func (a *AntiBruteForceServer) AddToBlackList(context.Context, *api.AddToBlackListRequest) (*api.AddToBlackListResponse, error) {
	panic("implement me")
}

func (a *AntiBruteForceServer) DeleteFromWhiteList(context.Context, *api.DeleteFromWhiteListRequest) (*api.DeleteFromWhiteListResponse, error) {
	panic("implement me")
}

func (a *AntiBruteForceServer) DeleteFromBlackList(context.Context, *api.DeleteFromBlackListRequest) (*api.DeleteFromBlackListResponse, error) {
	panic("implement me")
}

func (a *AntiBruteForceServer) ResetLimit(context.Context, *api.ResetLimitRequest) (*api.ResetLimitResponse, error) {
	panic("implement me")
}

func NewAntiBruteForceServer(antiBruteForceService *services.AntiBruteForceService) *AntiBruteForceServer {
	return &AntiBruteForceServer{AntiBruteForceService: antiBruteForceService}
}

func (a *AntiBruteForceServer) Serve(addr string) error {
	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGINT)
		<-stop
		log.Printf("Interrupt signal")
		log.Printf("Gracefully shutdown")
		s.GracefulStop()
	}()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	api.RegisterAntiBruteForceServiceServer(s, a)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s)
	return s.Serve(l)
}
