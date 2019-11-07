package models

import "net"

type Auth struct {
	Login    string
	Password string
	IpAddr   net.IP
}
