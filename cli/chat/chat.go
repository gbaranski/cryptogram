package chat

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Chat hold current chat state
type Chat struct {
	Context   *context.Context
	PubSub    *pubsub.PubSub
	Rooms     []*Room
	MsgSender *MessageSender
}

// CreateChat creates chat
func CreateChat(context context.Context, ps *pubsub.PubSub, defaultRoom *Room, nickname *string, peerID peer.ID) *Chat {
	var rooms []*Room
	rooms = append(rooms, defaultRoom)
	chat := &Chat{
		Context: &context,
		PubSub:  ps,
		Rooms:   rooms,
		MsgSender: &MessageSender{
			PeerID:   peerID,
			Nickname: *nickname,
		},
	}
	return chat

}
