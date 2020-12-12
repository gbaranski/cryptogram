package node

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// API used for holding current node state
type API struct {
	Host   *host.Host
	PubSub *pubsub.PubSub
}

// ChatRoomBufSize is the number of incoming messages to buffer for each topic.
const ChatRoomBufSize = 128

// ChatMessage gets converted to/from JSON and sent in the body of pubsub messages.
type ChatMessage struct {
	Message    *string
	SenderID   *string
	SenderNick *string
}

// ChatRoom represents a subscription to a single PubSub topic. Messages
// can be published to the topic with ChatRoom.Publish, and received
// messages are pushed to the Messages channel.
type ChatRoom struct {
	// Messages is a channel of messages received from other peers in the chat room

	Messages chan *ChatMessage

	Context *context.Context
	pubsub  *pubsub.PubSub
	topic   *pubsub.Topic
	sub     *pubsub.Subscription

	RoomName *string
	PeerID   *peer.ID
	Nickname *string
}

func getTopicName(roomName string) string {
	return "cryptogram-room:" + roomName
}

// JoinChatRoom tries to subscribe to the PubSub topic for the room name, returning
// a ChatRoom on success.
func JoinChatRoom(ctx context.Context, pubsub *pubsub.PubSub, peerID peer.ID, nickname *string, roomName *string) (*ChatRoom, error) {
	// join the pubsub topic
	topic, err := pubsub.Join(getTopicName(*roomName))
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &ChatRoom{
		Context:  &ctx,
		pubsub:   pubsub,
		topic:    topic,
		sub:      sub,
		PeerID:   &peerID,
		Nickname: nickname,
		RoomName: roomName,
		Messages: make(chan *ChatMessage, ChatRoomBufSize),
	}

	// start reading messages from the subscription in a loop
	go cr.readLoop()
	return cr, nil
}

// Publish sends a message to the pubsub topic.
func (cr *ChatRoom) Publish(message string) error {

	prettySenderID := cr.PeerID.Pretty()

	m := ChatMessage{
		Message:    &message,
		SenderID:   &prettySenderID,
		SenderNick: cr.Nickname,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.topic.Publish(*cr.Context, msgBytes)
}

// ListPeers list peers on specific channel
func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.pubsub.ListPeers(getTopicName(*cr.RoomName))
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *ChatRoom) readLoop() {
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
		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		cr.Messages <- cm
	}
}
