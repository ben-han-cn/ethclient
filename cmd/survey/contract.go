package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"goldenteam/ethclient"
	"goldenteam/ethclient/cluster"
)

var (
	ErrEmptyKeyPath    = errors.New("key path has no valid keys")
	ErrNodeIsNotSelect = errors.New("select node first")
	ErrUnknownAnswer   = errors.New("answer is unknown")
)

type Answer uint8

const (
	A Answer = iota
	B
	C
	D
)

func (a Answer) String() string {
	switch a {
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	default:
		panic("unknown answer")
	}
}

func AnswerFromString(a string) (Answer, error) {
	switch strings.ToUpper(a) {
	case "A":
		return A, nil
	case "B":
		return B, nil
	case "C":
		return C, nil
	case "D":
		return D, nil
	default:
		return 0, ErrUnknownAnswer
	}
}

type ContractClient struct {
	contractAddress  common.Address
	contractReady    bool
	node             *cluster.Node
	account          *ethclient.Account
	voteAccountIndex int
	client           *ethclient.Client
	gethpath         *cluster.GethPath
	contractABI      abi.ABI
	keyGenerator     *ethclient.KeyGenerator
}

func newContractClient(conf *cluster.Config, contract string) (*ContractClient, error) {
	keyFilePath, err := filepath.Abs(conf.KeyStorePath)
	if err != nil {
		return nil, err
	}

	keyGenerator := ethclient.NewKeyGenerator(keyFilePath)
	addresses := keyGenerator.ListAddress()
	if len(addresses) == 0 {
		return nil, ErrEmptyKeyPath
	}

	account, err := keyGenerator.GetAccount(addresses[0], cluster.DefaultPasswd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("use account %s\n", addresses[0].Hex())

	gethpath, err := cluster.NewGethPath(conf)
	if err != nil {
		return nil, err
	}

	contractABI, _ := ethclient.ABIFromString(SurveyABI)
	c := &ContractClient{
		account:          account,
		gethpath:         gethpath,
		contractABI:      contractABI,
		voteAccountIndex: 0,
		keyGenerator:     keyGenerator,
	}

	c.SelectNode("signer1")
	if contract == "" {
		c.deployContract()
	} else {
		c.contractAddress = common.HexToAddress(contract)
	}

	c.waitForContractReady()
	return c, nil
}

func (c *ContractClient) deployContract() {
	address, err := c.client.Deploy(c.account, SurveyABI, SurveyBin, "most used any encrypted currency")
	if err != nil {
		fmt.Printf("deploy contract failed %s\n", err.Error())
		os.Exit(1)
	} else {
		c.contractAddress = address
		fmt.Printf("deploy contract with address: %s\n", address.Hex())
	}
}

func (c *ContractClient) SelectNode(name string) error {
	if c.node != nil && c.node.Name() == name {
		return nil
	}

	node, err := cluster.NodeFromString(name)
	if err != nil {
		return err
	}

	c.node = node
	c.client = node.Client(c.gethpath)
	return nil
}

func (c *ContractClient) CurrentNodeName() string {
	if c.node == nil {
		return ""
	} else {
		return c.node.Name()
	}
}

func (c *ContractClient) waitForContractReady() {
	fmt.Printf("wait for contract ready\n")
	for {
		code, err := c.client.CodeAt(context.Background(), c.contractAddress, nil)
		if err == nil && len(code) > 0 {
			break
		}
		<-time.After(5 * time.Second)
	}
	fmt.Printf("contract is ready\n")
}

func (c *ContractClient) Topic() (string, error) {
	input, _ := ethclient.PackFunctionCall(c.contractABI, "topic")
	output, err := c.client.LocalCall(c.contractAddress, c.account, input)
	if err != nil {
		return "", err
	} else {
		var topic string
		err := c.contractABI.Unpack(&topic, "topic", output)
		return topic, err
	}
}

func (c *ContractClient) Vote(a Answer) error {
	addresses := c.keyGenerator.ListAddress()
	if c.voteAccountIndex >= len(addresses) {
		c.voteAccountIndex = 0
	}

	account, err := c.keyGenerator.GetAccount(addresses[c.voteAccountIndex], cluster.DefaultPasswd)
	if err != nil {
		return err
	}
	c.voteAccountIndex += 1

	input, _ := ethclient.PackFunctionCall(c.contractABI, "vote", uint8(a))
	return c.client.OnlineCall(c.contractAddress, account, input)
}

func (c *ContractClient) VoteForAnswer(a Answer) (*big.Int, error) {
	input, _ := ethclient.PackFunctionCall(c.contractABI, "voteForAnswer", uint8(a))
	output, err := c.client.LocalCall(c.contractAddress, c.account, input)
	if err != nil {
		return nil, err
	} else {
		number := big.NewInt(0)
		err := c.contractABI.Unpack(&number, "voteForAnswer", output)
		return number, err
	}
}

func (c *ContractClient) MostVoted() (Answer, error) {
	input, _ := ethclient.PackFunctionCall(c.contractABI, "mostVoted")
	output, err := c.client.LocalCall(c.contractAddress, c.account, input)
	if err != nil {
		return 0, err
	} else {
		var a uint8
		err := c.contractABI.Unpack(&a, "mostVoted", output)
		return Answer(a), err
	}
}

func (c *ContractClient) Finalize() error {
	input, _ := ethclient.PackFunctionCall(c.contractABI, "finalize")
	return c.client.OnlineCall(c.contractAddress, c.account, input)
}
