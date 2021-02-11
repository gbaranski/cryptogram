package node

import (
	core "github.com/libp2p/go-libp2p-core"
	discovery "github.com/libp2p/go-libp2p-discovery"
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
	Host   *core.Host
	PubSub *pubsub.PubSub
	DHT    *DhtAPI
}
