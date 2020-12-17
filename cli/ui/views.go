package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createInput(label *string, ch chan *string) *tview.InputField {
	input := tview.NewInputField().
		SetLabel(*label + " > ").
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

		// send the line onto the input chan and reset the field text
		ch <- &line
		input.SetText("")
	})

	return input
}

func setupInputHistory(input *tview.InputField, msgView *tview.TextView) *[]*string {
	var history []*string
	hc := 0
	input.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Modifiers() == tcell.ModAlt {
			row, _ := msgView.GetScrollOffset()
			switch e.Key() {
			case tcell.KeyDown:
				msgView.ScrollTo(row+1, 0)
			case tcell.KeyUp:
				msgView.ScrollTo(row-1, 0)
			}
			return e
		}
		switch e.Key() {
		case tcell.KeyUp:
			if hc >= len(history) || hc < 0 {
				return e
			}
			hc++
			input.SetText(*history[len(history)-hc])
		case tcell.KeyDown:
			if hc <= 1 {
				hc = 0
				input.SetText("")
				return e
			}
			hc--
			input.SetText(*history[len(history)-hc])
		default:
			hc = 0
		}
		return e
	})
	return &history
}

func createMsgView(drawFn func() *tview.Application) *tview.TextView {
	msgView := tview.NewTextView()
	msgView.SetDynamicColors(true)
	msgView.SetBorder(true)

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgView.SetChangedFunc(func() {
		drawFn()
	})
	return msgView
}

func createPeersView() *tview.TextView {
	peersView := tview.NewTextView()
	peersView.SetBorder(true).SetTitle("Peers")
	return peersView
}

func createChatPanel(msgView *tview.TextView, peersView *tview.TextView) *tview.Flex {
	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgView, 0, 1, false).
		AddItem(peersView, 20, 1, false)

	return chatPanel
}

func createFlex(chatPanel *tview.Flex, input *tview.InputField) *tview.Flex {
	// flex is a vertical box with the chatPanel on top and the input field at the bottom.
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)
}

func (ui *UI) updateRoomTitle() {
	ui.msgView.SetTitle(fmt.Sprintf("Room: %s", ui.room.Topic.String()))
}
