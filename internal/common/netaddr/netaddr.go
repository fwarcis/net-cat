package netaddr

import (
	"fmt"
	"net/netip"
	"strconv"
)

type IPHost string

const LocalHost = IPHost("localhost")

func (h *IPHost) UnmarshalText(text []byte) error {
	if string(text) == "localhost" {
		*h = IPHost(text)
		return nil
	}

	host, err := netip.ParseAddr(string(text))
	*h = IPHost(host.String())
	return err
}

func (h *IPHost) MarshalText() ([]byte, error) {
	text := []byte(string(*h))
	err := h.UnmarshalText(text)
	return text, err
}

func (h IPHost) String() string {
	text, err := h.MarshalText()
	_ = err
	return string(text)
}

type Port uint16

func (p *Port) UnmarshalText(text []byte) error {
	num, err := strconv.Atoi(string(text))
	if err != nil || num < 0 || num > 65535 {
		return fmt.Errorf(
			"invalid port '%s': must be a decimal between 0 and 65535",
			text)
	}

	*p = Port(num)
	return nil
}

func (p Port) MarshalText() ([]byte, error) {
	text := strconv.Itoa(int(p))
	return []byte(text), nil
}

func (p Port) String() string {
	return strconv.Itoa(int(p))
}
