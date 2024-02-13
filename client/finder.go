package client

import (
	"context"
	"net"

	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
)

const (
	//_service  = "_ledy._tcp"
	_hostname = "ledy-server"
	_domain   = ".local"
)

type finderResult struct {
	IP   net.IP
	Port int
}

// Find service on all network interfaces.
func findService(ctx context.Context) (*finderResult, error) {
	// https://github.com/pion/mdns/blob/master/examples/query/main.go
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		return nil, err
	}
	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}
	server, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{})
	if err != nil {
		return nil, err
	}
	_, src, err := server.Query(ctx, _hostname+_domain)
	return &finderResult{
		IP:   net.ParseIP(src.String()),
		Port: 81,
	}, err
}
