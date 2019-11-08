package models

import (
	"net"
)

type NetList struct {
	Networks []*net.IPNet
}
