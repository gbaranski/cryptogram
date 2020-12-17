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

func (ui *UI) handleHistory(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyUp:
		if ui.hc >= len(ui.history) || ui.hc < 0 {
			return
		}
		ui.hc++
		ui.inputField.SetText(*ui.history[len(ui.history)-ui.hc])
	case tcell.KeyDown:
		if ui.hc <= 1 {
			ui.hc = 0
			ui.inputField.SetText("")
		}
		ui.hc--
		ui.inputField.SetText(*ui.history[len(ui.history)-ui.hc])
	default:
		ui.hc = 0
	}
}

func (ui *UI) handleKeyScroll(e *tcell.EventKey) {
	row, _ := ui.msgView.GetScrollOffset()
	switch e.Key() {
	case tcell.KeyDown:
		ui.msgView.ScrollTo(row+1, 0)
	case tcell.KeyUp:
		ui.msgView.ScrollTo(row-1, 0)
	}
}

func (ui *UI) setupInputCapture() {
	ui.inputField.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyUp || e.Key() == tcell.KeyDown {
			if e.Modifiers() == tcell.ModAlt {
				ui.handleKeyScroll(e)
			} else {
				ui.handleHistory(e)
			}
		} else if e.Key() == tcell.KeyTAB {
			ui.app.SetFocus(ui.msgView)
		}
		return e
	})
	ui.msgView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyTAB {
			// Change to ui.peersView if this one on bottom would be uncommented
			ui.app.SetFocus(ui.inputField)
		}
		return e
	})

	// Its currently disabled, not sure if scrolling peers is nessesary
	/*
		ui.peersView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
			if e.Key() == tcell.KeyTAB {
				ui.app.SetFocus(ui.inputField)
			}
			return e
		})
	*/
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
