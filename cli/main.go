package main

import (
	"context"
	"log"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gbaranski/cryptogram/cli/cli"
	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/multiformats/go-multiaddr"

	node "github.com/gbaranski/cryptogram/cli/node"
)

func main() {

	var bootstrapPeers []multiaddr.Multiaddr

	for _, s := range []string{
		"/ip4/192.168.1.100/tcp/4001/p2p/QmNaRGMqFkSNEra3SQRFGQBZmHRmhCTq1ytuNsjkYqAGpQ",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		bootstrapPeers = append(bootstrapPeers, ma)
	}

	config := &misc.Config{
		RendezvousString: "cryptogram-rendezvous",
		ListenAddresses:  nil,
		ProtocolID:       "/chat/1.0.0",
		MDNSDiscovery: &misc.MDNSDiscoveryConfig{
			Enabled:  false,
			Interval: time.Minute * 15,
		},
		DHTDiscovery: &misc.DHTDiscoveryConfig{
			BootstrapPeers: &bootstrapPeers,
			Enabled:        true,
		},
	}

	log.Println("-- Getting an LibP2P host running -- ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := node.CreateAPI(&ctx, config)
	if err != nil {
		log.Panicf("Failed creating host: %s\n", err)
	}

	p := prompt.New(
		func(str string) { cli.Executor(str, api) },
		cli.Completer,
		prompt.OptionTitle("cryptogram-cli"),
		prompt.OptionPrefix(">>> "))
	p.Run()

}
