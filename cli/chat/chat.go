package chat

import (
	"context"

	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Chat hold current chat state
type Chat struct {
	Context   *context.Context
	PubSub    *pubsub.PubSub
	host      *host.Host
	MsgSender *MessageSender
}

// CreateChat creates chat
func CreateChat(context context.Context, ps *pubsub.PubSub, config *misc.Config, host *host.Host) *Chat {
	chat := &Chat{
		Context: &context,
		PubSub:  ps,
		host:    host,
		MsgSender: &MessageSender{
			PeerID:   (*host).ID(),
			Nickname: *config.Nickname,
		},
	}
	return chat
}

// ListAllPeers returns slice of all peers connected
func (chat *Chat) ListAllPeers() []peer.ID {
	return (*chat.host).Network().Peers()
}
