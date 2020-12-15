package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/gbaranski/cryptogram/cli/node"
	"github.com/gbaranski/cryptogram/cli/ui"
)

func main() {
	config := misc.GetConfig()
	ui := ui.CreateUI(config)
	go ui.RunApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := node.CreateAPI(&ctx, config, ui)

	if err != nil {
		log.Panicln("Failed creating host: ", err)
	}

	room, err := chat.CreateRoom(context.Background(), api.PubSub, config.Room, (*api.Host).ID())
	if err != nil {
		log.Panicln("Failed creating room ", err)
	}
	newChat := chat.CreateChat(ctx, api.PubSub, config, api.Host)

	if err != nil {
		log.Panicln("Error when creating chat: ", err)
	}
	ui.Log(fmt.Sprintf("Hi %s, use /help to get info about commands", *config.Nickname))
	ui.StartChat(newChat, room)
	<-ui.DoneCh
}
