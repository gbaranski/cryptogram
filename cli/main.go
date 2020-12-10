package main

import (
	"context"
	"log"
	"time"

	"github.com/gbaranski/cryptogram/cli/misc"
	dht "github.com/libp2p/go-libp2p-kad-dht"

	node "github.com/gbaranski/cryptogram/cli/node"
)

func main() {
	config := &misc.Config{
		RendezvousString: "cryptogram-rendezvous",
		BootstrapPeers:   dht.DefaultBootstrapPeers,
		ListenAddresses:  nil,
		ProtocolID:       "/chat/1.0.0",
		MDNSDiscovery: misc.MDNSDiscoveryConfig{
			Enabled:  true,
			Interval: time.Hour,
		},
		DHTDiscovery: misc.DHTDiscoveryConfig{
			Enabled: false,
		},
	}

	log.Println("-- Getting an LibP2P host running -- ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := node.CreateHost(&ctx, config)
	if err != nil {
		log.Panicf("Failed creating host: %s\n", err)
	}

	// p := prompt.New(func(str string) { repl.Executor(str, *ipfs) }, repl.Completer, prompt.OptionTitle("cryptogram-cli"), prompt.OptionPrefix(">>> "))
	// p.Run()
	select {}

}
