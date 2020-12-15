package ui

import (
	"log"

	chat "github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/rivo/tview"
)

// UI is a Text User Interface (TUI) for a ChatRoom.
// The Run method will draw the UI to the terminal in "fullscreen"
type UI struct {
	app       *tview.Application
	peersView *tview.TextView
	msgView   *tview.TextView
	history   *[]*string
	inputCh   chan *string
	doneCh    chan struct{}

	chat *chat.Chat
	room *chat.Room
}

// CreateUI returns UI
// It won't actually do anything until you call Run().
func CreateUI(config *misc.Config) *UI {
	app := tview.NewApplication()
	// an input field for typing messages into
	inputCh := make(chan *string, 32)
	input := createInput(config.Nickname, inputCh)
	history := setupInputHistory(input)
	msgView := createMsgView(app.Draw)
	peersView := createPeersView()
	chatPanel := createChatPanel(msgView, peersView)
	flex := createFlex(chatPanel, input)
	app.SetRoot(flex, true)

	return &UI{
		app:       app,
		peersView: peersView,
		history:   history,
		msgView:   msgView,
		inputCh:   inputCh,
		doneCh:    make(chan struct{}, 1),
	}
}

// RunApp starts UI app
func (ui *UI) RunApp() {
	defer ui.end()
	err := ui.app.Run()
	if err != nil {
		log.Panicln("Error when starting UI", err)
	}
}

// StartChat starts chat
func (ui *UI) StartChat(chat *chat.Chat, room *chat.Room) {
	ui.chat = chat
	ui.room = room
	ui.updateRoomTitle()
	go ui.handleEvents()
}

// end signals the event loop to exit gracefully
func (ui *UI) end() {
	ui.doneCh <- struct{}{}
}
