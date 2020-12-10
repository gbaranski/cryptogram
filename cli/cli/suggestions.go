package cli

import "github.com/c-bata/go-prompt"

var (
	commandSuggestions = []prompt.Suggest{
		{
			Text:        "peers",
			Description: "Show connected peers with PubSub",
		},
		{
			Text:        "id",
			Description: "Print current node ID",
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
)
