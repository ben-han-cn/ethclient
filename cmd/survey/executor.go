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
		case "topic":
			cmdTopic(c)
		case "vote":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("vote answer\n")
				return
			}
			cmdVote(c, cmdAndArgs[1])
		case "vote_for_answer":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("vote_for_answer answer\n")
				return
			}
			cmdVoteForAnswer(c, cmdAndArgs[1])
		case "most_voted":
			cmdMostVoted(c)
		case "finalize":
			cmdFinalize(c)
		case "quit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Printf("unknown cmd %s\n", cmd)
		}

		p.SetPrefix(c.CurrentNodeName() + ">>> ")
	}
}

func cmdTopic(c *ContractClient) {
	if topic, err := c.Topic(); err != nil {
		fmt.Printf("increment call failed:%s", err.Error())
	} else {
		fmt.Printf("%s\n", topic)
	}
}

func cmdVote(c *ContractClient, a string) {
	answer, err := AnswerFromString(a)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}

	err = c.Vote(answer)
	if err == nil {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("vote failed:%s\n", err.Error())
	}
}

func cmdVoteForAnswer(c *ContractClient, a string) {
	answer, err := AnswerFromString(a)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}

	num, err := c.VoteForAnswer(answer)
	if err == nil {
		fmt.Printf("%v\n", num)
	} else {
		fmt.Printf("err:%s\n", err.Error())
	}
}

func cmdMostVoted(c *ContractClient) {
	answer, err := c.MostVoted()
	if err == nil {
		fmt.Printf("%v\n", answer)
	} else {
		fmt.Printf("err:%s\n", err.Error())
	}
}

func cmdFinalize(c *ContractClient) {
	err := c.Finalize()
	if err == nil {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("err:%s\n", err.Error())
	}
}
