package node

import (
	"context"
	"fmt"

	"github.com/ipfs/go-ipfs/core"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

// IPFSNodeAPI Combines IPFS API and Node
type IPFSNodeAPI struct {
	API  iface.CoreAPI
	Node *core.IpfsNode
}

func getMessageTopic(receiverID string) string {
	return receiverID + "/messages"
}

// SendMessage sends a message via PubSub
func SendMessage(ipfs IPFSNodeAPI, topic string, msg string) {
	fmt.Printf("Sending message from %s with content: %s\n", topic, msg)
	// topic := getMessageTopic(receiverID)
	err := ipfs.API.PubSub().Publish(context.Background(), topic, []byte(msg))
	if err != nil {
		fmt.Printf("Error occured when publishing message %s\n", err)
	}
}

// SubscribeMessages subscribes to messages over PubSub
func SubscribeMessages(ipfs IPFSNodeAPI, receiverID string) iface.PubSubSubscription {
	subscription, err := ipfs.API.PubSub().Subscribe(context.TODO(), receiverID+"/messages")
	if err != nil {
		fmt.Printf("Error occured when subscribing %s\n", err)
	}
	return subscription
}
