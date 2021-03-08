package discovery

import (
	"errors"
	"net"
)

type MDNSConfig struct {
	Enabled bool
}

type MDNSClient struct {
	conn *net.UDPConn
}

func NewMDNSClient() (client MDNSClient, err error) {
	if client.conn, err = net.ListenMulticastUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   net.IPv4zero,
			Port: 5353,
			Zone: "",
		},
	); err != nil {
		return
	}
	if client.conn == nil {
		return client, errors.New("failed creating connection")
	}

	return
}
