package chat

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// RoomBufSize is the number of incoming messages to buffer for each topic.
const RoomBufSize = 128

// Message gets converted to/from JSON and sent in the body of pubsub messages.
type Message struct {
	Message    string
	SenderID   string
	SenderNick string
}

// Room represents a subscription to a single PubSub topic. Messages
// can be published to the topic with ChatRoom.Publish, and received
// messages are pushed to the Messages channel.
type Room struct {
	// Messages is a channel of messages received from other peers in the chat room

	Messages chan *Message

	Context *context.Context
	pubsub  *pubsub.PubSub
	topic   *pubsub.Topic
	sub     *pubsub.Subscription

	RoomName *string
	PeerID   *peer.ID
	Nickname *string
}

func getTopicName(roomName *string) string {
	return "cryptogram-room:" + *roomName
}

// JoinRoom tries to subscribe to the PubSub topic for the room name, returning
// a ChatRoom on success.
func JoinRoom(ctx context.Context, pubsub *pubsub.PubSub, peerID peer.ID, nickname *string, roomName string) (*Room, error) {
	// join the pubsub topic
	topic, err := pubsub.Join(getTopicName(&roomName))
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &Room{
		Context:  &ctx,
		pubsub:   pubsub,
		topic:    topic,
		sub:      sub,
		PeerID:   &peerID,
		Nickname: nickname,
		RoomName: &roomName,
		Messages: make(chan *Message, RoomBufSize),
	}

	// start reading messages from the subscription in a loop
	go cr.readLoop()
	return cr, nil
}

// Publish sends a message to the pubsub topic.
func (cr *Room) Publish(message string) error {

	m := Message{
		Message:    message,
		SenderID:   cr.PeerID.Pretty(),
		SenderNick: *cr.Nickname,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.topic.Publish(*cr.Context, msgBytes)
}

// ListPeers list peers on specific channel
func (cr *Room) ListPeers() []peer.ID {
	return cr.pubsub.ListPeers(getTopicName(cr.RoomName))
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *Room) readLoop() {
	for {
		msg, err := cr.sub.Next(*cr.Context)
		if err != nil {
			close(cr.Messages)
			return
		}
		// only forward messages delivered by others
		if &msg.ReceivedFrom == cr.PeerID {
			continue
		}
		cm := new(Message)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		cr.Messages <- cm
	}
}
