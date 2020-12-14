package main

import (
	"context"
	"log"

	"github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/misc"

	node "github.com/gbaranski/cryptogram/cli/node"
)

func main() {
	config := misc.GetConfig()
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
	newChat := chat.CreateChat(ctx, api.PubSub, room, config.Nickname, (*api.Host).ID())

	if err != nil {
		log.Panicln("Error when creating chat: ", err)
	}

	ui := chat.NewUI(newChat, room)
	if err = ui.Run(); err != nil {
		log.Panicln("error running text UI: ", err)
	}

}
