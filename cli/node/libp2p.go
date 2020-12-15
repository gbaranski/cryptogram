package node

import (
	"context"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/gbaranski/cryptogram/cli/discovery"
	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/gbaranski/cryptogram/cli/ui"
)

// API used for holding current node state
type API struct {
	Host   *host.Host
	PubSub *pubsub.PubSub
}

// CreateAPI creates libp2p API
func CreateAPI(ctx *context.Context, config *misc.Config, ui *ui.UI) (*API, error) {
	var opts []libp2p.Option
	opts = append(opts, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"))
	if *config.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	}
	host, err := libp2p.New(*ctx, opts...)
	if err != nil {
		return nil, err
	}
	ui.Log("LibP2P host is running ID: ", host.ID())

	if *config.Debug {
		ui.LogDebug("Host addresses: ")
		for _, addr := range host.Addrs() {
			ui.LogDebug(addr)
		}
	}
	ps, err := pubsub.NewGossipSub(*ctx, host)

	if config.DHTDiscovery.Enabled {
		ui.Log("Initializing DHT Discovery")
		_, _, err := discovery.SetupDHTDiscovery(ctx, &host, config, ui)
		if err != nil {
			return nil, err
		}
	}
	if config.MDNSDiscovery.Enabled {
		ui.Log("Initializing MDNS Discovery")
		err := discovery.SetupMDNSDiscovery(ctx, &host, config, ui)
		if err != nil {
			return nil, err
		}
	}

	return &API{Host: &host, PubSub: ps}, nil

}
