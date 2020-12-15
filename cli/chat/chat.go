package chat

import (
	"context"

	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Chat hold current chat state
type Chat struct {
	Context   *context.Context
	PubSub    *pubsub.PubSub
	Host      *host.Host
	MsgSender *MessageSender
}

// CreateChat creates chat
func CreateChat(context context.Context, ps *pubsub.PubSub, config *misc.Config, host *host.Host) *Chat {
	chat := &Chat{
		Context: &context,
		PubSub:  ps,
		Host:    host,
		MsgSender: &MessageSender{
			PeerID:   (*host).ID(),
			Nickname: *config.Nickname,
		},
	}
	return chat
}
