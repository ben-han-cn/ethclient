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
	{Text: "increment", Description: "Increment the number"},
	{Text: "get_number", Description: "Get current number"},
	{Text: "select_node", Description: "Call from which node"},
	{Text: "quit", Description: "quite the app"},
}
