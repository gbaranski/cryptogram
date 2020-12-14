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
	closeCh     chan struct{}
	context     *context.Context
	topic       *pubsub.Topic
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
	eh, err := topic.EventHandler()
	if err != nil {
		return nil, err
	}
	closeCh := make(chan struct{}, 1)
	room := &Room{
		msgCh:       make(chan *Message, RoomBufSize),
		peerEventCh: make(chan (*pubsub.PeerEvent)),
		closeCh:     closeCh,
		context:     &context,
		topic:       topic,
	}
	go room.readMessagesLoop(&peerID, subscription)
	go room.readEventsLoop(eh)
	go func() {
		select {
		case <-closeCh:
			eh.Cancel()
			subscription.Cancel()
			closeCh <- struct{}{}
			// close(room.msgCh)
			// close(room.peerEventCh)
			return
		}
	}()
	return room, nil
}

func (room *Room) readEventsLoop(eh *pubsub.TopicEventHandler) {
	for {
		e, err := eh.NextPeerEvent(*room.context)
		fmt.Println("Next peer event")
		if err != nil {
			return
		}
		room.peerEventCh <- &e
	}
}

func (room *Room) readMessagesLoop(peerID *peer.ID, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(*room.context)
		if err != nil {
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
	room.closeCh <- struct{}{}
	<-room.closeCh
	return room.topic.Close()
}
