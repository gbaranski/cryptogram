package misc

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Pallinder/go-randomdata"

	dht "github.com/libp2p/go-libp2p-kad-dht"
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
	Nickname         *string
	ListenAddresses  *[]maddr.Multiaddr
	ProtocolID       string
	MDNSDiscovery    *MDNSDiscoveryConfig
	DHTDiscovery     *DHTDiscoveryConfig
}

// GetConfig retreives config
// Can exit the program
func GetConfig() *Config {
	nickname := flag.String("nick", "Generated in runtime", "Nickname")
	dhtEnabled := flag.Bool("dht", false, "True if DHT discovery should be enabled")
	mdnsEnabled := flag.Bool("mdns", false, "True if MDNS discovery should be enabled")
	flag.Parse()

	if *dhtEnabled == false && *mdnsEnabled == false {
		fmt.Println("No discovery mode is enabled")
		fmt.Println("Please run with `-dht` or/and `-mdns` option")
		os.Exit(1)
	}
	// fix that stupid code later
	if *nickname == "Generated in runtime" {
		*nickname = fmt.Sprintf("%s-%s", randomdata.FirstName(randomdata.RandomGender), randomdata.LastName())
	}

	config := &Config{
		RendezvousString: "cryptogram-rendezvous",
		Nickname:         nickname,
		ListenAddresses:  nil,
		ProtocolID:       "/chat/1.0.0",
		MDNSDiscovery: &MDNSDiscoveryConfig{
			Enabled:  *mdnsEnabled,
			Interval: time.Minute * 15,
		},
		DHTDiscovery: &DHTDiscoveryConfig{
			BootstrapPeers: &dht.DefaultBootstrapPeers,
			Enabled:        *dhtEnabled,
		},
	}
	return config

}
