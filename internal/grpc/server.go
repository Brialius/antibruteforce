package grpc

import (
	"context"
	"github.com/Brialius/antibruteforce/internal/domain/errors"
	"github.com/Brialius/antibruteforce/internal/domain/models"
	"github.com/Brialius/antibruteforce/internal/domain/services"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ip := net.ParseIP(req.GetAuth().GetIp())
	if ip == nil {
		log.Printf("IP address is invalid: %s", req.GetAuth().GetIp())
		return &api.CheckAuthResponse{
			Result: &api.CheckAuthResponse_Error{
				Error: errors.ErrInvalidIP.Error(),
			},
		}, nil
	}
	ok, err := a.AntiBruteForceService.CheckAuth(ctx, &models.Auth{
		Login:    req.GetAuth().GetLogin(),
		Password: req.GetAuth().GetPassword(),
		IpAddr:   ip,
	})
	if err != nil {
		log.Printf("Error during checking auth `%s`: %s", req.GetAuth(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.CheckAuthResponse{
		Result: &api.CheckAuthResponse_Ok{Ok: ok},
	}, nil
}

func (a *AntiBruteForceServer) AddToWhiteList(ctx context.Context, req *api.AddToWhiteListRequest) (*api.AddToWhiteListResponse, error) {
	_, n, err := net.ParseCIDR(req.GetNet())
	if err != nil {
		log.Printf("Error during adding `%s` to whitelist: %s", req.GetNet(), err)
		return &api.AddToWhiteListResponse{
			Error: errors.ErrInvalidCIDR.Error(),
		}, nil
	}
	err = a.AntiBruteForceService.AddToWhiteList(ctx, n)
	if err != nil {
		log.Printf("Error during adding `%s` to whitelist: %s", req.GetNet(), err)
		if berr, ok := err.(errors.AntiBruteForceError); ok {
			resp := &api.AddToWhiteListResponse{
				Error: string(berr),
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.AddToWhiteListResponse{}, nil
}

func (a *AntiBruteForceServer) AddToBlackList(ctx context.Context, req *api.AddToBlackListRequest) (*api.AddToBlackListResponse, error) {
	_, n, err := net.ParseCIDR(req.GetNet())
	if err != nil {
		log.Printf("Error during adding `%s` to blacklist: %s", req.GetNet(), err)
		return &api.AddToBlackListResponse{
			Error: errors.ErrInvalidCIDR.Error(),
		}, nil
	}
	err = a.AntiBruteForceService.AddToBlackList(ctx, n)
	if err != nil {
		log.Printf("Error during adding `%s` to blacklist: %s", req.GetNet(), err)
		if berr, ok := err.(errors.AntiBruteForceError); ok {
			resp := &api.AddToBlackListResponse{
				Error: string(berr),
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.AddToBlackListResponse{}, nil
}

func (a *AntiBruteForceServer) DeleteFromWhiteList(ctx context.Context, req *api.DeleteFromWhiteListRequest) (*api.DeleteFromWhiteListResponse, error) {
	_, n, err := net.ParseCIDR(req.GetNet())
	if err != nil {
		log.Printf("Error during deleting `%s` from whitelist: %s", req.GetNet(), err)
		return &api.DeleteFromWhiteListResponse{
			Error: errors.ErrInvalidCIDR.Error(),
		}, nil
	}
	err = a.AntiBruteForceService.DeleteFromWhiteList(ctx, n)
	if err != nil {
		log.Printf("Error during deleting `%s` from whitelist: %s", req.GetNet(), err)
		if berr, ok := err.(errors.AntiBruteForceError); ok {
			resp := &api.DeleteFromWhiteListResponse{
				Error: string(berr),
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.DeleteFromWhiteListResponse{}, nil
}

func (a *AntiBruteForceServer) DeleteFromBlackList(ctx context.Context, req *api.DeleteFromBlackListRequest) (*api.DeleteFromBlackListResponse, error) {
	_, n, err := net.ParseCIDR(req.GetNet())
	if err != nil {
		log.Printf("Error during deleting `%s` from blacklist: %s", req.GetNet(), err)
		return &api.DeleteFromBlackListResponse{
			Error: errors.ErrInvalidCIDR.Error(),
		}, nil
	}
	err = a.AntiBruteForceService.DeleteFromBlackList(ctx, n)
	if err != nil {
		log.Printf("Error during deleting `%s` from blacklist: %s", req.GetNet(), err)
		if berr, ok := err.(errors.AntiBruteForceError); ok {
			resp := &api.DeleteFromBlackListResponse{
				Error: string(berr),
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.DeleteFromBlackListResponse{}, nil
}

func (a *AntiBruteForceServer) ResetLimit(ctx context.Context, req *api.ResetLimitRequest) (*api.ResetLimitResponse, error) {
	ip := net.ParseIP(req.GetIp())
	if ip == nil {
		log.Printf("Error during resetting limit for ip `%s`: Invalid IP", req.GetIp())
		return &api.ResetLimitResponse{
			Error: errors.ErrInvalidCIDR.Error(),
		}, nil
	}

	err := a.AntiBruteForceService.ResetLimit(ctx, req.GetLogin(), &ip)
	if err != nil {
		log.Printf("Error during resetting limit for login `%s` and IP `%s`: %s", req.GetLogin(), req.GetIp(), err)
		if berr, ok := err.(errors.AntiBruteForceError); ok {
			resp := &api.ResetLimitResponse{
				Error: string(berr),
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.ResetLimitResponse{}, nil
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
