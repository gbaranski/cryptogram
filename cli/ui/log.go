package ui

import (
	"fmt"

	chat "github.com/gbaranski/cryptogram/cli/chat"
	"github.com/gbaranski/cryptogram/cli/misc"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Log logs to UI
func (ui *UI) Log(args ...interface{}) {
	prompt := misc.WithColor("red", "<System>:")
	fmt.Fprintln(ui.msgView, prompt, fmt.Sprint(args...))
}

// LogError logs error
func (ui *UI) LogError(when string, err error) {
	prompt := misc.WithColor("darkred", "<Error>:")
	fmt.Fprintln(ui.msgView, fmt.Sprintf("%s Error occured when %s: %s", prompt, when, err.Error()))
}

// LogDebug similar to Log but logs only if debug mode
func (ui *UI) LogDebug(args ...interface{}) {
	if !*ui.config.Debug {
		return
	}
	prompt := misc.WithColor("blue", "<Debug>:")
	fmt.Fprintln(ui.msgView, prompt, fmt.Sprint(args...))
}

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
func (ui *UI) logChatMessage(cm *chat.Message) {
	var color string
	if cm.Sender.PeerID.Pretty() == ui.chat.MsgSender.PeerID.Pretty() {
		color = "yellow"
	} else {
		color = "green"
	}
	if len(cm.Sender.Nickname) < 1 {
		cm.Sender.Nickname = "Unknown"
	}
	prompt := misc.WithColor(color, fmt.Sprintf("<%s>:", cm.Sender.Nickname))
	fmt.Fprintf(ui.msgView, "%s %s\n", prompt, cm.Text)
}

func (ui *UI) logPeerEvent(e *pubsub.PeerEvent) {
	switch e.Type {
	case pubsub.PeerJoin:
		ui.Log(fmt.Sprintf("%s joined", e.Peer.Pretty()))
	case pubsub.PeerLeave:
		ui.Log(fmt.Sprintf("%s left", e.Peer.Pretty()))
	}
}
