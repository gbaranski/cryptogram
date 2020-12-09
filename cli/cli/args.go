package cli

import (
	"flag"

	misc "github.com/gbaranski/cryptogram/cli/misc"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// GetConfig retreives config from args
func GetConfig() (misc.Config, error) {
	config := misc.Config{}

	flag.StringVar(&config.RendezvousString, "rendezvous", "meet me here",
		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.StringVar(&config.ProtocolID, "pid", "/chat/1.0.0", "Sets a protocol id for stream headers")
	flag.Parse()

	if len(config.BootstrapPeers) == 0 {
		config.BootstrapPeers = dht.DefaultBootstrapPeers
	}

	return config, nil

}
