package repl

import "github.com/c-bata/go-prompt"

var (
	commandSuggestions = []prompt.Suggest{
		{
			Text:        "send_message",
			Description: "Send messages",
		},
		{
			Text:        "topic",
			Description: "Manage topics",
		},
		{
			Text:        "help",
			Description: "Show available commands",
		},
		{
			Text:        "exit",
			Description: "Exit the repl",
		},
	}

	topicSuggestions = []prompt.Suggest{
		{
			Text:        "subscribe",
			Description: "Subscribe to a topic",
		},
		{
			Text:        "unsubscribe",
			Description: "Unsubscribe from a topic",
		},
	}
)
