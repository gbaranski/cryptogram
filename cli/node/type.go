package node

import (
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// DhtAPI is API for DHT things
type DhtAPI struct {
	Discovery *discovery.RoutingDiscovery
	IpfsDHT   *dht.IpfsDHT
}

// API used for holding current node state
type API struct {
	Host   *host.Host
	PubSub *pubsub.PubSub
	DHT    *DhtAPI
}
