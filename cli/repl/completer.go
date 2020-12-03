package repl

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// Completer adds completer feature
func Completer(d prompt.Document) []prompt.Suggest {
	if strings.HasPrefix(d.Text, "topic") {
		return topicSuggestions
	}
	return commandSuggestions
}
