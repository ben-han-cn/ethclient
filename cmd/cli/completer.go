package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"goldenteam/ethclient/cluster"
)

func NewCompleter(ctrl *cluster.Controller) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		args := strings.Split(d.TextBeforeCursor(), " ")
		if len(args) == 1 {
			results := prompt.FilterHasPrefix(commands, d.GetWordBeforeCursor(), true)
			return results
		}

		cmd := args[0]
		arg := d.GetWordBeforeCursor()
		return completeCmdArgs(ctrl, cmd, arg)
	}
}

var commands = []prompt.Suggest{
	{Text: "nodes", Description: "Display all the nodes"},
	{Text: "stopnode", Description: "stop the node"},
	{Text: "restartnode", Description: "restart the stopped node"},
	{Text: "addsyncer", Description: "add new syncer"},
	{Text: "addsigner", Description: "add new signer"},
	{Text: "selectnode", Description: "Connect to one node"},
	{Text: "blocknumber", Description: "get current block number"},
	{Text: "block", Description: "get block"},
	{Text: "balance", Description: "Get balance of one account"},
	{Text: "transaction", Description: "Get transaction with hash"},
	{Text: "transfer", Description: "Transfer money between random accounts"},
	{Text: "quit", Description: "quite the app"},
}

func completeCmdArgs(ctrl *cluster.Controller, cmd, arg string) []prompt.Suggest {
	switch cmd {
	case "selectnode", "stopnode":
		nodes := append(ctrl.Signers().Nodes(), ctrl.Syncers().Nodes()...)
		nodesSuggest := make([]prompt.Suggest, len(nodes))
		for i, node := range nodes {
			nodesSuggest[i].Text = node.Name()
		}
		return prompt.FilterHasPrefix(nodesSuggest, arg, true)
	case "restartnode":
		nodes := append(ctrl.Signers().MissingNodes(), ctrl.Syncers().MissingNodes()...)
		nodesSuggest := make([]prompt.Suggest, len(nodes))
		for i, node := range nodes {
			nodesSuggest[i].Text = node.Name()
		}
		return prompt.FilterHasPrefix(nodesSuggest, arg, true)
	case "balance":
		accounts := ctrl.Accounts()
		accountSuggest := make([]prompt.Suggest, len(accounts))
		for i, account := range accounts {
			accountSuggest[i].Text = account.Hex()
		}
		return prompt.FilterHasPrefix(accountSuggest, arg, true)
	default:
		return nil
	}
}
