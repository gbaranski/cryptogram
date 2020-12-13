package chat

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI is a Text User Interface (TUI) for a ChatRoom.
// The Run method will draw the UI to the terminal in "fullscreen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type UI struct {
	chat      *Chat
	room      *Room
	app       *tview.Application
	peersList *tview.TextView
	msgW      io.Writer
	inputCh   chan string
	doneCh    chan struct{}
}

// NewUI returns a new ChatUI struct that controls the text UI.
// It won't actually do anything until you call Run().
func NewUI(chat *Chat, room *Room) *UI {
	app := tview.NewApplication()

	msgView := tview.NewTextView()
	msgView.SetDynamicColors(true)
	msgView.SetBorder(true)
	msgView.SetTitle(fmt.Sprintf("Room: %s", room.topic.String()))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgView.SetChangedFunc(func() {
		app.Draw()
	})

	// an input field for typing messages into
	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(chat.msgSender.Nickname + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		// send the line onto the input chan and reset the field text
		inputCh <- line
		input.SetText("")
	})

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgView, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)

	return &UI{
		chat:      chat,
		room:      room,
		app:       app,
		peersList: peersList,
		msgW:      msgView,
		inputCh:   inputCh,
		doneCh:    make(chan struct{}, 1),
	}
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *UI) Run() error {
	go ui.handleEvents()
	defer ui.end()

	return ui.app.Run()
}

// end signals the event loop to exit gracefully
func (ui *UI) end() {
	ui.doneCh <- struct{}{}
}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *UI) refreshPeers() {
	peers := ui.room.topic.ListPeers()
	idStrs := make([]string, len(peers))
	for i, p := range peers {
		idStrs[i] = p.String()
	}

	ui.peersList.SetText(strings.Join(idStrs, "\n"))
	ui.app.Draw()
}

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
func (ui *UI) displayChatMessage(cm *Message) {
	var color string
	if cm.Sender.PeerID.Pretty() == ui.chat.msgSender.PeerID.Pretty() {
		color = "yellow"
	} else {
		color = "green"
	}
	if len(cm.Sender.Nickname) < 1 {
		cm.Sender.Nickname = "Unknown"
	}
	prompt := withColor(color, fmt.Sprintf("<%s>:", cm.Sender.Nickname))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, cm.Text)
}

func (ui *UI) handleCommand(command string) {
	args := strings.Split(command, " ")
	msg := &Message{
		Sender: MessageSender{
			Nickname: "Cryptogram-bot",
		},
	}
	switch args[0] {
	case "join":
		ui.room.close()
		room, err := CreateRoom(context.Background(), ui.chat.pubsub, args[1], ui.chat.msgSender.PeerID)
		if err != nil {
			log.Panicln("Error when creating new room ", err)
		}
		ui.room = room
		ui.refreshPeers()
	case "help":
		msg.Text = "List of commands\n/join <room-name>"
	default:
		msg.Text = "Unknown command, use /help to list commands"
	}

	ui.room.msgChan <- msg

}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.
func (ui *UI) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	for {
		select {
		case input := <-ui.inputCh:
			if strings.HasPrefix(input, "/") {
				ui.handleCommand(strings.TrimPrefix(input, "/"))
				continue
			}
			// when the user types in a line, publish it to the chat room and print to the message window
			message := &Message{
				Text:   input,
				Sender: *ui.chat.msgSender,
			}
			err := ui.room.sendMessage(context.Background(), message)
			if err != nil {
				log.Panicf("publish error: %s\n", err)
			}
		case m := <-ui.room.msgChan:
			// when we receive a message from the chat room, print it to the message window
			ui.displayChatMessage(m)

		case <-peerRefreshTicker.C:
			// refresh the list of peers in the chat room periodically
			ui.refreshPeers()
		case <-(*ui.room.context).Done():
			fmt.Println("Context done")
			return

		case <-ui.doneCh:
			return
		}
	}
}

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
