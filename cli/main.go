package main

import (
	"context"
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
	ui.Log("-- Getting an LibP2P host running -- ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := node.CreateAPI(&ctx, config, ui)

	if err != nil {
		log.Panicln("Failed creating host: ", err)
	}

	room, err := chat.CreateRoom(context.Background(), api.PubSub, "general", (*api.Host).ID())
	if err != nil {
		log.Panicln("Failed creating room ", err)
	}
	newChat := chat.CreateChat(ctx, api.PubSub, room, config.Nickname, (*api.Host).ID())

	if err != nil {
		log.Panicln("Error when creating chat: ", err)
	}
	ui.StartChat(newChat, room)
	select {}
}
