package chat

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Room is room
type Room struct {
	msgChan      chan *Message
	doneChan     chan struct{}
	context      *context.Context
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

// CreateRoom creates Room for desired RoomName
func CreateRoom(context context.Context, pubsub *pubsub.PubSub, roomName string, peerID peer.ID) (*Room, error) {
	topic, err := pubsub.Join(GetTopicName(&roomName))
	if err != nil {
		return nil, err
	}
	subscription, err := topic.Subscribe()
	room := &Room{
		msgChan:      make(chan *Message, RoomBufSize),
		context:      &context,
		topic:        topic,
		subscription: subscription,
	}
	go room.readLoop(&peerID)
	return room, nil

}

func (room *Room) readLoop(peerID *peer.ID) {
	for {
		msg, err := room.subscription.Next(*room.context)
		if err != nil {
			close(room.msgChan)
			return
		}

		// only forward messages delivered by others
		// if msg.ReceivedFrom == *peerID {
		// 	continue
		// }
		cm := new(Message)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		room.msgChan <- cm
	}
}

func (room *Room) sendMessage(context context.Context, message *Message) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = room.topic.Publish(context, msgBytes)
	return err
}
