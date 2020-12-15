package discovery

import (
	"context"
	"sync"
	"time"

	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/gbaranski/cryptogram/cli/ui"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func connectBootstrapNodes(ctx *context.Context, host *host.Host, config *misc.Config, ui *ui.UI) error {
	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range *config.DHTDiscovery.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := (*host).Connect(*ctx, *peerinfo); err != nil {
				ui.LogError("connecting to bootstrap node", err)
			} else {
				ui.Log("Connection established with bootstrap node: ", *peerinfo)
			}
		}()
	}
	wg.Wait()
	return nil

}

func searchPeers(
	ctx *context.Context,
	host *host.Host,
	routingDiscovery *discovery.RoutingDiscovery,
	config *misc.Config,
	ui *ui.UI) {
	ui.LogDebug("Searching for peers")
	peerChan, err := routingDiscovery.FindPeers(*ctx, *config.RendezvousName)
	if err != nil {
		panic(err)
	}
	for p := range peerChan {
		if p.ID == (*host).ID() {
			continue
		}

		ui.LogDebug("DHT Peer found ID: ", p.ID)
		go (*host).Connect(context.Background(), p)
	}
}

func searchPeersLoop(
	ctx *context.Context,
	host *host.Host,
	routingDiscovery *discovery.RoutingDiscovery,
	config *misc.Config,
	ui *ui.UI) {
	// Do the init search
	searchPeers(ctx, host, routingDiscovery, config, ui)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				searchPeers(ctx, host, routingDiscovery, config, ui)
			}
		}
	}()
}

// SetupDHTDiscovery set ups DHT Discovery
func SetupDHTDiscovery(ctx *context.Context, host *host.Host, config *misc.Config, ui *ui.UI) (*discovery.RoutingDiscovery, *dht.IpfsDHT, error) {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(*ctx, *host)
	if err != nil {
		return nil, nil, err
	}
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	if err = kademliaDHT.Bootstrap(*ctx); err != nil {
		return nil, nil, err
	}

	err = connectBootstrapNodes(ctx, host, config, ui)
	if err != nil {
		return nil, nil, err
	}

	// We use a rendezvous point "cryptogram-rendezvous" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(*ctx, routingDiscovery, *config.RendezvousName)
	ui.LogDebug("Successfully announced ourselfs, other peers can now find us!")

	go searchPeers(ctx, host, routingDiscovery, config, ui)

	return routingDiscovery, kademliaDHT, nil
}
