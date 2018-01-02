package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"goldenteam/ethclient"
	"goldenteam/ethclient/cluster"
)

var (
	ErrEmptyKeyPath    = errors.New("key path has no valid keys")
	ErrNodeIsNotSelect = errors.New("select node first")
)

type ContractClient struct {
	contractAddress common.Address
	contractReady   bool
	node            *cluster.Node
	account         *ethclient.Account
	client          *ethclient.Client
	gethpath        *cluster.GethPath
	contractABI     abi.ABI
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

	contractABI, _ := ethclient.ABIFromString(IncrementerABI)
	c := &ContractClient{
		account:     account,
		gethpath:    gethpath,
		contractABI: contractABI,
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
	address, err := c.client.Deploy(c.account, IncrementerABI, IncrementerBin)
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

func (c *ContractClient) Increment() error {
	input, _ := ethclient.PackFunctionCall(c.contractABI, "increment")
	return c.client.OnlineCall(c.contractAddress, c.account, input)
}

func (c *ContractClient) GetNumber() (*big.Int, error) {
	input, _ := ethclient.PackFunctionCall(c.contractABI, "getNumber")
	output, err := c.client.LocalCall(c.contractAddress, c.account, input)
	if err != nil {
		return nil, err
	} else {
		number := big.NewInt(0)
		err := c.contractABI.Unpack(&number, "getNumber", output)
		return number, err
	}
}
