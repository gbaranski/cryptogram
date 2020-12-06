package repl

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type mathSuggestion struct {
	promptSuggestion prompt.Suggest
	score            int
}

func findSuggestions(str string, availableSuggestions []prompt.Suggest) []prompt.Suggest {
	var suggestionsStrings []string
	for _, e := range availableSuggestions {
		suggestionsStrings = append(suggestionsStrings, e.Text)
	}
	foundSuggestionsString := fuzzy.Find(str, suggestionsStrings)
	var foundSuggestions []prompt.Suggest
	for _, e := range foundSuggestionsString {
		for _, c := range availableSuggestions {
			if e == c.Text {
				foundSuggestions = append(foundSuggestions, c)
			}

		}

	}
	return foundSuggestions
}

// Completer adds completer feature
func Completer(d prompt.Document) []prompt.Suggest {
	cmdText := strings.Trim(d.Text, " ")
	if strings.HasPrefix(cmdText, "topic") {
		cmdText := strings.Trim(strings.TrimPrefix(cmdText, "topic"), " ")
		return findSuggestions(cmdText, topicSuggestions)
	}
	return findSuggestions(cmdText, commandSuggestions)
}
