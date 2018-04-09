package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/ethereum/go-ethereum/common"
	"goldenteam/ethclient"
	"goldenteam/ethclient/cluster"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewExecutor(ctrl *cluster.Controller) func(*prompt.Prompt, string) {
	var currentNode *cluster.Node
	return func(p *prompt.Prompt, in string) {
		in = strings.TrimSpace(in)
		if in == "" {
			return
		}

		cmdAndArgs := strings.Split(in, " ")
		cmd := cmdAndArgs[0]
		switch cmd {
		case "nodes":
			cmdListNodes(ctrl)
		case "addsyncer":
			cmdAddSyncer(ctrl)
		case "addsigner":
			cmdAddSigner(ctrl)
		case "stopnode":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("stopnode node\n")
				return
			}
			n := cmdStopNode(ctrl, cmdAndArgs[1])
			if n != nil && currentNode != nil && currentNode.Equals(n) {
				currentNode = nil
			}
		case "restartnode":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("restartnode node\n")
				return
			}
			cmdRestartNode(ctrl, cmdAndArgs[1])
		case "selectnode":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("connect node\n")
				return
			}
			if n := cmdSelectNode(ctrl, cmdAndArgs[1]); n != nil {
				currentNode = n
			}
		case "blocknumber":
			cmdBlockNumber(ctrl, currentNode)
		case "block":
			if len(cmdAndArgs) == 1 {
				cmdBlock(ctrl, currentNode, -1)
			} else {
				number, err := strconv.Atoi(cmdAndArgs[1])
				if err != nil {
					fmt.Printf("block number isn't valid int\n")
					return
				}
				cmdBlock(ctrl, currentNode, number)
			}
		case "transaction":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("transaction hash\n")
				return
			}
			cmdTransaction(ctrl, currentNode, cmdAndArgs[1])

		case "balance":
			if len(cmdAndArgs) != 2 {
				fmt.Printf("balance account\n")
				return
			}
			cmdBalance(ctrl, currentNode, cmdAndArgs[1])
		case "transfer":
			if len(cmdAndArgs) < 2 {
				fmt.Printf("transfer amount of money\n")
				return
			}

			value, err := strconv.Atoi(cmdAndArgs[1])
			if err != nil {
				fmt.Printf("amount isn't valid number\n")
				return
			}

			count := 1
			if len(cmdAndArgs) == 3 {
				count, err = strconv.Atoi(cmdAndArgs[2])
			}
			cmdTransfer(ctrl, int64(value), count)
		case "quit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Printf("unknown cmd %s\n", cmd)
		}

		if currentNode == nil {
			p.SetPrefix(">>> ")
		} else {
			p.SetPrefix(currentNode.Name() + ">>> ")
		}

	}
}

func cmdListNodes(ctrl *cluster.Controller) {
	for _, n := range ctrl.Signers().Nodes() {
		fmt.Printf("%s\n", n.Name())
	}

	fmt.Printf("+++++++\n")

	for _, n := range ctrl.Syncers().Nodes() {
		fmt.Printf("%s\n", n.Name())
	}
}

func cmdAddSyncer(ctrl *cluster.Controller) {
	if _, err := ctrl.AddSyncer(); err == nil {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("err:%s\n", err.Error())
	}
}

func cmdAddSigner(ctrl *cluster.Controller) {
	if _, err := ctrl.AddSigner(); err == nil {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("err:%s\n", err.Error())
	}
}

func cmdStopNode(ctrl *cluster.Controller, name string) *cluster.Node {
	if n, err := cluster.NodeFromString(name); err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return nil
	} else {
		if err := ctrl.StopNode(n); err != nil {
			fmt.Printf("err:%s\n", err.Error())
			return nil
		} else {
			fmt.Printf("ok\n")
			return n
		}
	}
}

func cmdRestartNode(ctrl *cluster.Controller, name string) {
	if n, err := cluster.NodeFromString(name); err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		if err := ctrl.RestartNode(n); err != nil {
			fmt.Printf("err:%s\n", err.Error())
		} else {
			fmt.Printf("ok\n")
		}
	}
}

func cmdSelectNode(ctrl *cluster.Controller, name string) *cluster.Node {
	if n, err := cluster.NodeFromString(name); err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return nil
	} else {
		if _, err := ctrl.ClientForNode(n); err != nil {
			fmt.Printf("err:%s\n", err.Error())
			return nil
		} else {
			fmt.Printf("ok\n")
			return n
		}
	}
}

func cmdBlock(ctrl *cluster.Controller, node *cluster.Node, number int) {
	client := getClient(ctrl, node)
	if client == nil {
		return
	}

	var blockNumber *big.Int
	if number >= 0 {
		blockNumber = big.NewInt(int64(number))
	}

	block, err := client.BlockByNumber(blockNumber)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		blockMarshal := cluster.BlockInPOA(block, ctrl)
		fmt.Printf(`block:         %v
signer:        %s
hash:          %s
parent:        %s
transactions:  %v
difficulty:    %v
time:          %v
is_vote:       %v
voite_address: %v
is_epoch:      %v
signers:       %v
`, blockMarshal.Number,
			blockMarshal.Signer,
			blockMarshal.Hash,
			blockMarshal.Parent,
			blockMarshal.Transactions,
			blockMarshal.Difficulty,
			blockMarshal.Time,
			blockMarshal.IsVote,
			blockMarshal.VoteAddress,
			blockMarshal.IsEpoch,
			blockMarshal.Signers)
	}
}

func cmdBlockNumber(ctrl *cluster.Controller, node *cluster.Node) {
	client := getClient(ctrl, node)
	if client == nil {
		return
	}

	if number, err := client.BlockNumber(); err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		fmt.Printf("%v\n", number)
	}
}

func cmdTransaction(ctrl *cluster.Controller, node *cluster.Node, hash string) {
	client := getClient(ctrl, node)
	if client == nil {
		return
	}

	tx, _, err := client.TransactionByHash(common.HexToHash(hash))
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	} else {
		txMarshal := cluster.TransactionInPOA(tx)
		fmt.Printf(`tx:        %s
from:      %s
to:        %s
value:     %v
data:      %x
`, txMarshal.Hash,
			txMarshal.From,
			txMarshal.To,
			txMarshal.Value,
			txMarshal.Data)
	}

	r, err := client.TransactionReceipt(common.HexToHash(hash))
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		receiptMarshal := cluster.ReceiptInPOA(r)
		fmt.Printf(`succeed:   %v
logs:      %v
`, receiptMarshal.Succeed,
			receiptMarshal.Logs)
	}
}

func cmdTransfer(ctrl *cluster.Controller, value int64, count int) {
	accounts := ctrl.Accounts()
	fromIndex := rand.Intn(len(accounts))
	toIndex := rand.Intn(len(accounts))
	for toIndex == fromIndex {
		toIndex = rand.Intn(len(accounts))
	}
	from := accounts[fromIndex]
	to := accounts[toIndex]

	if err := ctrl.TransferMoney(from, to, value, count); err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("done transfer from %s to %s\n", from.Hex(), to.Hex())
	}
}

func cmdBalance(ctrl *cluster.Controller, node *cluster.Node, address string) {
	client := getClient(ctrl, node)
	if client == nil {
		return
	}

	if balance, err := client.GetBalance(common.HexToAddress(address)); err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		fmt.Printf("%v\n", balance)
	}
}

func getClient(ctrl *cluster.Controller, node *cluster.Node) *ethclient.Client {
	if node == nil {
		fmt.Printf("select node first\n")
		return nil
	}

	client, err := ctrl.ClientForNode(node)
	if err != nil {
		fmt.Printf("node isn't valid\n")
	}
	return client
}
