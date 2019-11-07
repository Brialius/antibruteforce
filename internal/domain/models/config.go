package models

import "net"

type NetList struct {
	Networks []*net.IPNet
}

type WhiteList struct {
	NetList
}

type BlackList struct {
	NetList
}
