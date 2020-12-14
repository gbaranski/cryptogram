package discovery

import (
	"context"
	"log"
	"sync"

	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func connectBootstrapNodes(ctx *context.Context, host *host.Host, config *misc.Config) error {

	log.Println("\n-- Going to connect to a few nodes in the Network as bootstrappers --")
	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range *config.DHTDiscovery.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := (*host).Connect(*ctx, *peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()
	return nil

}

func findAndConnectPeers(
	ctx *context.Context,
	host *host.Host,
	routingDiscovery *discovery.RoutingDiscovery,
	config *misc.Config) {
	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	// log.Println("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(*ctx, *config.RendezvousString)
	if err != nil {
		panic(err)
	}
	for p := range peerChan {
		if p.ID == (*host).ID() {
			continue
		}
		go (*host).Connect(context.Background(), p)
	}

}

// SetupDHTDiscovery set ups DHT Discovery
func SetupDHTDiscovery(ctx *context.Context, host *host.Host, config *misc.Config) (*discovery.RoutingDiscovery, *dht.IpfsDHT, error) {
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
	log.Println("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(*ctx); err != nil {
		return nil, nil, err
	}

	err = connectBootstrapNodes(ctx, host, config)
	if err != nil {
		return nil, nil, err
	}

	// We use a rendezvous point "cryptogram-rendezvous" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	log.Println("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(*ctx, routingDiscovery, *config.RendezvousString)
	log.Println("Successfully announced!")

	go findAndConnectPeers(ctx, host, routingDiscovery, config)

	return routingDiscovery, kademliaDHT, nil
}
