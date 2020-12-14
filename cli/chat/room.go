package chat

import (
	"context"
	"encoding/json"
	"log"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Room is room
type Room struct {
	msgCh        chan *Message
	peerEventCh  chan *pubsub.PeerEvent
	context      *context.Context
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
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
	room := &Room{
		msgCh:        make(chan *Message, RoomBufSize),
		peerEventCh:  make(chan (*pubsub.PeerEvent)),
		context:      &context,
		topic:        topic,
		subscription: subscription,
	}
	go room.readMessagesLoop(&peerID)
	go room.readEventsLoop(eh)
	return room, nil

}

func (room *Room) readEventsLoop(eh *pubsub.TopicEventHandler) {
	for {
		e, err := eh.NextPeerEvent(*room.context)
		if err != nil {
			select {
			case <-room.peerEventCh:
				log.Panicln("Event channel is already closed", err)
				return
			default:
				close(room.msgCh)
				return
			}
		}
		room.peerEventCh <- &e
	}

}

func (room *Room) readMessagesLoop(peerID *peer.ID) {
	for {

		msg, err := room.subscription.Next(*room.context)
		if err != nil {
			select {
			case <-room.msgCh:
				log.Panicln("Error when tracking sub", err)
				return
			default:
				close(room.msgCh)
				return
			}
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
	return err
}

func (room *Room) close() {
	room.subscription.Cancel()
	room.topic.Close()
}
