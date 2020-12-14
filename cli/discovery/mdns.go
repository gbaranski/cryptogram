package discovery

import (
	"context"
	"log"

	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type discoveryNotifee struct {
	Host *host.Host
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	err := (*(n.Host)).Connect(context.Background(), pi)
	if err != nil {
		log.Panicln("Error when connecting to MDNS peer", err)
		return
	}
}

// SetupMDNSDiscovery set ups MDNS Discovery
func SetupMDNSDiscovery(ctx *context.Context, host *host.Host, config *misc.Config) error {
	disc, err := discovery.NewMdnsService(*ctx, *host, config.MDNSDiscovery.Interval, *config.RendezvousString)

	if err != nil {
		return err
	}
	n := &discoveryNotifee{Host: host}
	disc.RegisterNotifee(n)
	return nil
}
