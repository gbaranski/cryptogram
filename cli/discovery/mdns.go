package discovery

import (
	"context"

	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// SetupMDNSDiscovery set ups MDNS Discovery
func SetupMDNSDiscovery(ctx *context.Context, host *host.Host, config *misc.Config) (chan peer.AddrInfo, error) {
	disc, err := discovery.NewMdnsService(*ctx, *host, config.MDNSDiscovery.Interval, config.RendezvousString)

	if err != nil {
		return nil, err
	}
	n := &discoveryNotifee{PeerChan: make(chan peer.AddrInfo)}
	disc.RegisterNotifee(n)
	return n.PeerChan, nil
}
