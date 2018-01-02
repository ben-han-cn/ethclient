package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func NewCompleter(c *ContractClient) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		args := strings.Split(d.TextBeforeCursor(), " ")
		if len(args) == 1 {
			results := prompt.FilterHasPrefix(commands, d.GetWordBeforeCursor(), true)
			return results
		} else {
			return nil
		}
	}
}

var commands = []prompt.Suggest{
	{Text: "vote", Description: "Vote for answer"},
	{Text: "vote_for_answer", Description: "Vote for answer"},
	{Text: "most_voted", Description: "Get most voted answer"},
	{Text: "finalize", Description: "Finish this round of suvery"},
	{Text: "topic", Description: "Get survey topic"},
	{Text: "select_node", Description: "Call from which node"},
	{Text: "quit", Description: "quite the app"},
}
