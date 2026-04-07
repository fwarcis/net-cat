package conf

import "netcat/internal/common/netaddr"

type config struct {
	Host netaddr.IPHost
	Port netaddr.Port
}

func NewDefaultConfig() config {
	return config{
		Host: netaddr.LocalHost,
		Port: netaddr.Port(8989),
	}
}
