package models

import "net"

// Auth domain model
type Auth struct {
	Login    string
	Password string
	IPAddr   net.IP
}
