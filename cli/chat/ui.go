package chat

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/gdamore/tcell/v2"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
	peersView *tview.TextView
	msgView   *tview.TextView
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
	peersView := tview.NewTextView()
	peersView.SetBorder(true)
	peersView.SetTitle("Peers")

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgView, 0, 1, false).
		AddItem(peersView, 20, 1, false)

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
		peersView: peersView,
		msgView:   msgView,
		inputCh:   inputCh,
		doneCh:    make(chan struct{}, 1),
	}
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *UI) Run() error {
	ui.updateRoomTitle()
	go ui.handleEvents()
	defer ui.end()

	return ui.app.Run()
}

func (ui *UI) updateRoomTitle() {
	ui.msgView.SetTitle(fmt.Sprintf("Room: %s", ui.room.topic.String()))
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

	ui.peersView.SetText(strings.Join(idStrs, "\n"))
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
	fmt.Fprintf(ui.msgView, "%s %s\n", prompt, cm.Text)
}

func (ui *UI) displayPeerEvent(e *pubsub.PeerEvent) {
	switch e.Type {
	case pubsub.PeerJoin:
		ui.printSystemMessage(fmt.Sprintf("%s joined", e.Peer.Pretty()))
	case pubsub.PeerLeave:
		ui.printSystemMessage(fmt.Sprintf("%s left", e.Peer.Pretty()))
	}
}

func (ui *UI) printSystemMessage(args ...interface{}) {
	prompt := withColor("red", "<System>:")
	fmt.Fprintln(ui.msgView, prompt, fmt.Sprint(args...))
}

func (ui *UI) printErrorMessage(when string, err error) {
	prompt := withColor("darkred", "<Error>:")
	fmt.Fprintln(ui.msgView, fmt.Sprintf("%s Error occured when %s: %s", prompt, when, err.Error()))
}

func (ui *UI) handleCommand(command string) {
	args := strings.Split(command, " ")
	switch args[0] {
	case "join":
		err := ui.room.close()
		if err != nil {
			ui.printErrorMessage("closing room", err)
			return
		}
		room, err := CreateRoom(context.Background(), ui.chat.pubsub, args[1], ui.chat.msgSender.PeerID)
		if err != nil {
			ui.printErrorMessage("creating room", err)
			return
		}
		ui.room = room
		ui.refreshPeers()
		ui.updateRoomTitle()
		ui.printSystemMessage("Successfully joined room: ", args[1])
	case "help":
		ui.printSystemMessage(`Commands:
		  /join <room-name>		- Joins a room
		  /topics 				 - Prints out all subscribed topics
		  /free					- Removes garbage from memory
		  /stats				 - Prints statistics
		  `)
	case "free":
		memStats := misc.GetMemStats()
		runtime.GC()
		ui.msgView.Clear()
		newMemStats := misc.GetMemStats()
		ui.printSystemMessage(fmt.Sprintf("Freed %fMiB of memory", memStats.Alloc-newMemStats.Alloc))
	case "stats":
		memStats := misc.GetMemStats()
		ui.printSystemMessage(fmt.Sprintf("Heap alloc - %fMiB", memStats.Alloc))
	case "topics":
		for i, t := range ui.chat.pubsub.GetTopics() {
			ui.printSystemMessage("Subscribed topics: ")
			fmt.Fprintf(ui.msgView, "%d - %s\n", i, t)
		}
	default:
		ui.printSystemMessage("Unknown command, use /help to list commands")
	}
}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.
func (ui *UI) handleEvents() {
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
				ui.printErrorMessage("sending message", err)
			}
		case m := <-ui.room.msgCh:
			// when we receive a message from the chat room, print it to the message window
			if m != nil {
				ui.displayChatMessage(m)
			}
		case e := <-ui.room.peerEventCh:
			if e != nil {
				ui.displayPeerEvent(e)
				ui.refreshPeers()
			}
		case <-(*ui.room.context).Done():
			fmt.Println("Context done")
			return

		case <-ui.doneCh:
			return
		}
	}
}

// WithColor wraps a string with color tags for display in the messages text box. Only for UI
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
