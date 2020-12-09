package node

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/libp2p/go-libp2p"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"

	misc "github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			log.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		log.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			log.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			log.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func handleStream(stream network.Stream) {
	log.Println("Got a new stream!")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go readData(rw)
	go writeData(rw)
}

func connectBootstrapNodes(ctx *context.Context, host *host.Host, config *misc.Config) error {

	log.Println("\n-- Going to connect to a few nodes in the Network as bootstrappers --")
	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
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
	log.Println("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(*ctx, config.RendezvousString)
	if err != nil {
		panic(err)
	}
	for peer := range peerChan {
		if peer.ID == (*host).ID() {
			continue
		}
		log.Println("Found peer:", peer)

		log.Println("Connecting to:", peer)
		stream, err := (*host).NewStream(*ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			log.Println("Connection failed:", err)
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}

		log.Println("Connected to:", peer)
	}

}

func setupDHT(ctx *context.Context, host *host.Host) (*dht.IpfsDHT, error) {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(*ctx, *host)
	if err != nil {
		return nil, err
	}
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	log.Println("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(*ctx); err != nil {
		return nil, err
	}

	return kademliaDHT, nil
}

func announce(ctx *context.Context, dht *dht.IpfsDHT, config *misc.Config) *discovery.RoutingDiscovery {
	// We use a rendezvous point "cryptogram-rendezvous" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	log.Println("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(dht)
	discovery.Advertise(*ctx, routingDiscovery, config.RendezvousString)
	log.Println("Successfully announced!")

	return routingDiscovery
}

// CreateHost creates libp2p host
func CreateHost(ctx *context.Context, config *misc.Config) (*host.Host, error) {
	host, err := libp2p.New(*ctx)
	if err != nil {
		return nil, err
	}
	log.Println("LibP2P host is running ID:", host.ID())

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)
	dht, err := setupDHT(ctx, &host)
	if err != nil {
		return nil, err
	}

	err = connectBootstrapNodes(ctx, &host, config)
	if err != nil {
		return nil, err
	}

	routingDiscovery := announce(ctx, dht, config)
	go findAndConnectPeers(ctx, &host, routingDiscovery, config)

	return &host, nil

}
