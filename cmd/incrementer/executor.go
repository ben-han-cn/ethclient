package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

func NewExecutor(c *ContractClient) func(*prompt.Prompt, string) {
	return func(p *prompt.Prompt, in string) {
		in = strings.TrimSpace(in)
		if in == "" {
			return
		}

		cmdAndArgs := strings.Split(in, " ")
		cmd := cmdAndArgs[0]
		switch cmd {
		case "increment":
			cmdIncrement(c)
		case "get_number":
			cmdGetNumber(c)
		case "select_node":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("select_node node_name\n")
				return
			}
			cmdSelectNode(c, cmdAndArgs[1])
		case "quit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Printf("unknown cmd %s\n", cmd)
		}

		p.SetPrefix(c.CurrentNodeName() + ">>> ")
	}
}

func cmdIncrement(c *ContractClient) {
	if err := c.Increment(); err != nil {
		fmt.Printf("increment call failed:%s", err.Error())
	} else {
		fmt.Printf("ok\n")
	}
}

func cmdGetNumber(c *ContractClient) {
	number, err := c.GetNumber()
	if err == nil {
		fmt.Printf("%v\n", number)
	} else {
		fmt.Printf("get number failed:%s\n", err.Error())
	}
}

func cmdSelectNode(c *ContractClient, name string) {
	if err := c.SelectNode(name); err == nil {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("err:%s\n", err.Error())
	}
}
