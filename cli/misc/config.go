package misc

import (
	"time"

	maddr "github.com/multiformats/go-multiaddr"
)

// MDNSDiscoveryConfig used for Config
type MDNSDiscoveryConfig struct {
	Enabled  bool
	Interval time.Duration
}

// DHTDiscoveryConfig used for Config
type DHTDiscoveryConfig struct {
	Enabled        bool
	BootstrapPeers *[]maddr.Multiaddr
}

// Config configuration
type Config struct {
	// Unique string to identify group of nodes. Share this with your friends to let them connect with you
	RendezvousString string
	ListenAddresses  *[]maddr.Multiaddr
	ProtocolID       string
	MDNSDiscovery    *MDNSDiscoveryConfig
	DHTDiscovery     *DHTDiscoveryConfig
}
