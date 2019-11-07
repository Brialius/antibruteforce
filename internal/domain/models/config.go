package models

import (
	"github.com/Brialius/antibruteforce/internal/domain/interfaces"
	"net"
)

type Config struct {
	WhiteList
	BlackList
	Storage *interfaces.ConfigStorage
}

type NetList struct {
	Networks []*net.IPNet
}

type WhiteList struct {
	NetList
}

type BlackList struct {
	NetList
}
