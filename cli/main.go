package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/discovery"
	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/gbaranski/cryptogram/cli/node"
	"github.com/gbaranski/cryptogram/cli/ui"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func main() {
	config := misc.GetConfig()
	ui := ui.CreateUI(config)
	go ui.RunApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := createAPI(&ctx, config, ui)

	if err != nil {
		log.Panicln("Failed creating host: ", err)
	}

	room, err := chat.CreateRoom(context.Background(), api.PubSub, config.Room, (*api.Host).ID())
	if err != nil {
		log.Panicln("Failed creating room ", err)
	}
	newChat := chat.CreateChat(ctx, config, api)

	if err != nil {
		log.Panicln("Error when creating chat: ", err)
	}
	ui.Log(fmt.Sprintf("Hi %s, use /help to get info about commands", *config.Nickname))
	ui.StartChat(newChat, room)
	<-ui.DoneCh
}

func createAPI(ctx *context.Context, config *misc.Config, ui *ui.UI) (*node.API, error) {
	var opts []libp2p.Option
	opts = append(opts, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"))
	if *config.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	}
	host, err := libp2p.New(*ctx, opts...)
	if err != nil {
		return nil, err
	}
	ui.Log("LibP2P host is running ID: ", host.ID())

	if *config.Debug {
		ui.LogDebug("Host addresses: ")
		for _, addr := range host.Addrs() {
			ui.LogDebug(addr)
		}
	}
	ps, err := pubsub.NewGossipSub(*ctx, host)
	var dhtAPI *node.DhtAPI

	if config.DHTDiscovery.Enabled {
		ui.Log("Initializing DHT Discovery")
		routingDiscovery, ipfsDHT, err := discovery.SetupDHTDiscovery(ctx, &host, config, ui)
		if err != nil {
			return nil, err
		}
		dhtAPI = &node.DhtAPI{Discovery: routingDiscovery, IpfsDHT: ipfsDHT}
	}
	if config.MDNSDiscovery.Enabled {
		ui.Log("Initializing MDNS Discovery")
		err := discovery.SetupMDNSDiscovery(ctx, &host, config, ui)
		if err != nil {
			return nil, err
		}
	}

	return &node.API{Host: &host, PubSub: ps, DHT: dhtAPI}, nil

}
