package interfaces

import (
	"context"
	"net"
)

type ConfigStorage interface {
	AddToBlackList(ctx context.Context, net *net.IPNet) error
	DeleteFromBlackList(ctx context.Context, net *net.IPNet) error
	AddToWhiteList(ctx context.Context, net *net.IPNet) error
	DeleteFromWhiteList(ctx context.Context, net *net.IPNet) error
	CheckIP(ctx context.Context, ip net.IP) bool
	Close(ctx context.Context)
}
