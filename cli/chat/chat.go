package chat

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Chat hold current chat state
type Chat struct {
	context   *context.Context
	pubsub    *pubsub.PubSub
	rooms     []*Room
	nickname  *string
	msgSender *MessageSender
}

// CreateChat creates chat
func CreateChat(context context.Context, ps *pubsub.PubSub, defaultRoom *Room, nickname *string, peerID peer.ID) *Chat {
	var rooms []*Room
	rooms = append(rooms, defaultRoom)
	chat := &Chat{
		context: &context,
		pubsub:  ps,
		rooms:   rooms,
		msgSender: &MessageSender{
			PeerID:   peerID,
			Nickname: *nickname,
		},
	}
	return chat

}
