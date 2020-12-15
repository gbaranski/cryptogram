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
	RendezvousName *string
	Nickname       *string
	Insecure       *bool
	MDNSDiscovery  *MDNSDiscoveryConfig
	DHTDiscovery   *DHTDiscoveryConfig
	Room           *string
}

// GetConfig retreives config
// Can exit the program
func GetConfig() *Config {
	usrEnv, usrEnvExists := os.LookupEnv("USER")
	var nickname *string
	if usrEnvExists {
		nickname = flag.String("nick", usrEnv, "Nickname to identify yourself")
	} else {
		nickname = flag.String("nick", "Randomly generated", "Nickname to identify yourself")
	}

	dhtEnabled := flag.Bool("dht", false, "True if DHT discovery should be enabled")
	mdnsEnabled := flag.Bool("mdns", false, "True if MDNS discovery should be enabled")
	mdnsInterval := flag.Int("mdnsinterval", 900, "MDNS Interval in seconds, by default 15 minutes")
	rendezvousName := flag.String("rendezvousname", "cryptogram-rendezvous", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	insecure := flag.Bool("insecure", false, "Disables security transports, incompatible with other secure transports, use with caution")
	room := flag.String("room", "general", "Default room you join at the start, can be changed later using /join command")
	flag.Parse()

	if *dhtEnabled == false && *mdnsEnabled == false {
		fmt.Println("No discovery mode is enabled")
		fmt.Println("Please run with `-dht` or/and `-mdns` option")
		os.Exit(1)
	}
	if !usrEnvExists {
		randomNickname := fmt.Sprintf("%s-%s", randomdata.FirstName(randomdata.RandomGender), randomdata.LastName())
		nickname = &randomNickname
	}

	config := &Config{
		RendezvousName: rendezvousName,
		Nickname:       nickname,
		Insecure:       insecure,
		Room:           room,
		MDNSDiscovery: &MDNSDiscoveryConfig{
			Enabled:  *mdnsEnabled,
			Interval: time.Duration(*mdnsInterval) * time.Second,
		},
		DHTDiscovery: &DHTDiscoveryConfig{
			BootstrapPeers: &dht.DefaultBootstrapPeers,
			Enabled:        *dhtEnabled,
		},
	}
	return config

}
