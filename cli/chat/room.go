package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Room is room
type Room struct {
	MsgCh       chan *Message
	PeerEventCh chan *pubsub.PeerEvent
	Context     *context.Context
	Topic       *pubsub.Topic
	Sub         *pubsub.Subscription
	Teh         *pubsub.TopicEventHandler
}

// CreateRoom creates Room for desired RoomName
func CreateRoom(context context.Context, ps *pubsub.PubSub, roomName *string, peerID peer.ID) (*Room, error) {
	topic, err := ps.Join(GetTopicName(roomName))
	if err != nil {
		return nil, err
	}
	subscription, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}
	teh, err := topic.EventHandler()
	if err != nil {
		return nil, err
	}
	room := &Room{
		MsgCh:       make(chan *Message, RoomBufSize),
		PeerEventCh: make(chan (*pubsub.PeerEvent)),
		Context:     &context,
		Topic:       topic,
		Sub:         subscription,
		Teh:         teh,
	}
	go room.readMessagesLoop(&peerID)
	go room.readPeerEventsLoop()
	return room, nil
}

func (room *Room) readPeerEventsLoop() {
	for {
		e, err := room.Teh.NextPeerEvent(*room.Context)
		if err != nil {
			fmt.Println("Error occured readPeerEventsLoop", err)
			return
		}
		room.PeerEventCh <- &e
	}
}

func (room *Room) readMessagesLoop(peerID *peer.ID) {
	for {
		msg, err := room.Sub.Next(*room.Context)
		if err != nil {
			// fmt.Println("Error occured readMessagesLoop", err)
			return
		}
		cm := new(Message)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			log.Panicln("Failed unmarshalling message", err)
		}
		// send valid messages onto the Messages channel
		room.MsgCh <- cm
	}
}

// SendMessage sends message
func (room *Room) SendMessage(context context.Context, message *Message) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = room.Topic.Publish(context, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

// Close closes room
func (room *Room) Close() error {
	room.Sub.Cancel()
	room.Teh.Cancel()
	return room.Topic.Close()
}
