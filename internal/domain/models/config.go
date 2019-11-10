package models

import (
	"net"
)

type NetList struct {
	Networks map[string]*net.IPNet
}
