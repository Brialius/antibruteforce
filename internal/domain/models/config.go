package models

import (
	"net"
)

// NetList model for white/black lists
type NetList struct {
	Networks map[string]*net.IPNet
}
