package interfaces

import (
	"context"
	"net"
)

type ConfigStorage interface {
	AddToBlackList(ctx context.Context, net net.IPNet)
	DeleteFromBlackList(ctx context.Context, net net.IPNet)
	AddToWhiteList(ctx context.Context, net net.IPNet)
	DeleteFromWhiteList(ctx context.Context, net net.IPNet)
	Close(ctx context.Context)
}
