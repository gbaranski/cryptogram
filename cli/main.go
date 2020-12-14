package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/misc"
	dht "github.com/libp2p/go-libp2p-kad-dht"

	node "github.com/gbaranski/cryptogram/cli/node"
)

func main() {
	nick := flag.String("nick", "Anonymous", "Your nickname")
	flag.Parse()

	config := &misc.Config{
		RendezvousString: "cryptogram-rendezvous",
		ListenAddresses:  nil,
		ProtocolID:       "/chat/1.0.0",
		MDNSDiscovery: &misc.MDNSDiscoveryConfig{
			Enabled:  true,
			Interval: time.Minute * 15,
		},
		DHTDiscovery: &misc.DHTDiscoveryConfig{
			BootstrapPeers: &dht.DefaultBootstrapPeers,
			Enabled:        false,
		},
	}

	log.Println("-- Getting an LibP2P host running -- ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := node.CreateAPI(&ctx, config)
	if err != nil {
		log.Panicln("Failed creating host: ", err)
	}

	room, err := chat.CreateRoom(context.Background(), api.PubSub, "general", (*api.Host).ID())
	if err != nil {
		log.Panicln("Failed creating room ", err)
	}
	newChat := chat.CreateChat(ctx, api.PubSub, room, nick, (*api.Host).ID())

	if err != nil {
		log.Panicln("Error when creating chat: ", err)
	}

	ui := chat.NewUI(newChat, room)
	if err = ui.Run(); err != nil {
		log.Panicln("error running text UI: ", err)
	}

}
