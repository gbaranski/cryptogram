package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gbaranski/cryptogram/cli/cli"
	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/multiformats/go-multiaddr"

	node "github.com/gbaranski/cryptogram/cli/node"
)

func main() {

	var bootstrapPeers []multiaddr.Multiaddr
	for _, s := range []string{
		"/ip4/192.168.1.100/tcp/4001/p2p/QmS4MVeG7LmTjW3NemtPLMsEmDKuTPcsx4Uk9H6EbAhtSV",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		bootstrapPeers = append(bootstrapPeers, ma)
	}

	config := &misc.Config{
		RendezvousString: "cryptogram-rendezvous",
		ListenAddresses:  nil,
		ProtocolID:       "/chat/1.0.0",
		MDNSDiscovery: &misc.MDNSDiscoveryConfig{
			Enabled:  true,
			Interval: time.Minute * 15,
		},
		DHTDiscovery: &misc.DHTDiscoveryConfig{
			BootstrapPeers: &bootstrapPeers,
			Enabled:        false,
		},
	}

	log.Println("-- Getting an LibP2P host running -- ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := node.CreateAPI(&ctx, config)
	if err != nil {
		log.Panicf("Failed creating host: %s\n", err)
	}
	rand.Seed(time.Now().Unix())
	nick := strconv.Itoa(rand.Int())
	log.Println("Nick!", nick)
	roomName := "hello-world"

	cr, err := node.JoinChatRoom(ctx, api.PubSub, (*api.Host).ID(), &nick, &roomName)
	if err != nil {
		log.Panicln("Error when joining chat room: ", err)
	}

	ui := cli.NewChatUI(cr)
	if err = ui.Run(); err != nil {
		log.Fatalf("error running text UI: %s", err)
	}

}
