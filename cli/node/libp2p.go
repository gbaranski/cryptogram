package node

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/gbaranski/cryptogram/cli/discovery"
	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// CreateAPI creates libp2p API
func CreateAPI(ctx *context.Context, config *misc.Config) (*API, error) {
	host, err := libp2p.New(*ctx)
	if err != nil {
		return nil, err
	}
	log.Println("LibP2P host is running ID:", host.ID())
	log.Println("Host addresses: ")
	for _, addr := range host.Addrs() {
		log.Println(addr)
	}
	ps, err := pubsub.NewGossipSub(*ctx, host)

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(config.ProtocolID), misc.HandleNetworkStream)
	if config.DHTDiscovery.Enabled {
		log.Println("Initializing DHT Discovery")
		_, _, err := discovery.SetupDHTDiscovery(ctx, &host, config)
		if err != nil {
			return nil, err
		}
	}
	if config.MDNSDiscovery.Enabled {
		log.Println("Initializing MDNS Discovery")
		peerChan, err := discovery.SetupMDNSDiscovery(ctx, &host, config)

		if err != nil {
			return nil, err
		}
		log.Println("Waiting for first peer to connect")
		peer := <-peerChan
		stream, err := host.NewStream(*ctx, peer.ID, protocol.ID(config.ProtocolID))
		if err != nil {
			return nil, err
		}
		misc.HandleNetworkStream(stream)

	}

	return &API{Host: &host, PubSub: ps}, nil

}
