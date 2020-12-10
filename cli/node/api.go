package node

import (
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// API used for holding current node state
type API struct {
	Host   *host.Host
	PubSub *pubsub.PubSub
}
