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
	Host     *host.Host
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Println("MDNS Peer found: ", pi.ID)
	err := (*(n.Host)).Connect(context.Background(), pi)
	if err != nil {
		log.Panicln("Error when connecting to MDNS peer", err)
		return
	}
	log.Println("Connected to: ", pi.ID)
	n.PeerChan <- pi
}

// SetupMDNSDiscovery set ups MDNS Discovery
func SetupMDNSDiscovery(ctx *context.Context, host *host.Host, config *misc.Config) (chan peer.AddrInfo, error) {
	disc, err := discovery.NewMdnsService(*ctx, *host, config.MDNSDiscovery.Interval, config.RendezvousString)

	if err != nil {
		return nil, err
	}
	n := &discoveryNotifee{PeerChan: make(chan peer.AddrInfo), Host: host}
	disc.RegisterNotifee(n)
	return n.PeerChan, nil
}
