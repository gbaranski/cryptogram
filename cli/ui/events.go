package ui

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	chat "github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/misc"
)

func (ui *UI) handleCommand(command string) {
	args := strings.Split(command, " ")
	switch args[0] {
	case "join":
		err := ui.room.Close()
		if err != nil {
			ui.LogError("closing room", err)
			return
		}
		room, err := chat.CreateRoom(context.Background(), ui.chat.PubSub, args[1], ui.chat.MsgSender.PeerID)
		if err != nil {
			ui.LogError("creating room", err)
			return
		}
		ui.room = room
		ui.refreshPeers()
		ui.updateRoomTitle()
		ui.Log("Successfully joined room: ", args[1])
	case "help":
		ui.Log(`Commands:
		  /join <room-name>		- Joins a room
		  /topics 				 - Prints out all subscribed topics
		  /free					- Removes garbage from memory
		  /stats				   - Prints memory statistics
		  `)
	case "free":
		memStats := misc.GetMemStats()
		runtime.GC()
		ui.msgView.Clear()
		newMemStats := misc.GetMemStats()
		ui.Log(fmt.Sprintf("Freed %fMiB of memory", memStats.Alloc-newMemStats.Alloc))
	case "stats":
		memStats := misc.GetMemStats()
		ui.Log(fmt.Sprintf("Heap alloc - %fMiB", memStats.Alloc))
	case "topics":
		for i, t := range ui.chat.PubSub.GetTopics() {
			ui.Log("Subscribed topics: ")
			fmt.Fprintf(ui.msgView, "%d - %s\n", i, t)
		}
	case "exit":
		ui.end()
	default:
		ui.Log("Unknown command, use /help to list commands")
	}
}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.
func (ui *UI) handleEvents() {
	for {
		select {
		case input := <-ui.inputCh:
			// Append to history only if the last str in history isn't the same as input
			if len(*ui.history) == 0 || *(*ui.history)[len(*ui.history)-1] != *input {
				*ui.history = append(*ui.history, input)
			}
			if strings.HasPrefix(*input, "/") {
				ui.handleCommand(strings.TrimPrefix(*input, "/"))
				continue
			}
			// when the user types in a line, publish it to the chat room and print to the message window
			message := &chat.Message{
				Text:   *input,
				Sender: *ui.chat.MsgSender,
			}
			err := ui.room.SendMessage(context.Background(), message)
			if err != nil {
				ui.LogError("sending message", err)
			}
		case m := <-ui.room.MsgCh:
			// when we receive a message from the chat room, print it to the message window
			if m != nil {
				ui.logChatMessage(m)
			}
		case e := <-ui.room.PeerEventCh:
			if e != nil {
				ui.logPeerEvent(e)
				ui.refreshPeers()
			}
		case <-ui.DoneCh:
			return
		}
	}
}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *UI) refreshPeers() {
	peers := ui.room.Topic.ListPeers()
	idStrs := make([]string, len(peers))
	for i, p := range peers {
		idStrs[i] = p.String()
	}

	ui.peersView.SetText(strings.Join(idStrs, "\n"))
	ui.app.Draw()
}
