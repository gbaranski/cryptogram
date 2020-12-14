package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Room is room
type Room struct {
	msgCh       chan *Message
	peerEventCh chan *pubsub.PeerEvent
	context     *context.Context
	topic       *pubsub.Topic
	sub         *pubsub.Subscription
	teh         *pubsub.TopicEventHandler
}

// CreateRoom creates Room for desired RoomName
func CreateRoom(context context.Context, ps *pubsub.PubSub, roomName string, peerID peer.ID) (*Room, error) {
	topic, err := ps.Join(GetTopicName(&roomName))
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
		msgCh:       make(chan *Message, RoomBufSize),
		peerEventCh: make(chan (*pubsub.PeerEvent)),
		context:     &context,
		topic:       topic,
		sub:         subscription,
		teh:         teh,
	}
	go room.readMessagesLoop(&peerID)
	go room.readPeerEventsLoop()
	return room, nil
}

func (room *Room) readPeerEventsLoop() {
	for {
		e, err := room.teh.NextPeerEvent(*room.context)
		if err != nil {
			fmt.Println("Error occured readPeerEventsLoop", err)
			return
		}
		room.peerEventCh <- &e
	}
}

func (room *Room) readMessagesLoop(peerID *peer.ID) {
	for {
		msg, err := room.sub.Next(*room.context)
		if err != nil {
			// fmt.Println("Error occured readMessagesLoop", err)
			return
		}
		cm := new(Message)
		err = json.Unmarshal(msg.Data, cm)
		// send valid messages onto the Messages channel
		room.msgCh <- cm
	}
}

func (room *Room) sendMessage(context context.Context, message *Message) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = room.topic.Publish(context, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (room *Room) close() error {
	room.sub.Cancel()
	room.teh.Cancel()
	return room.topic.Close()
}
