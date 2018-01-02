package ethclient

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) OnlineCall(contract common.Address, from *Account, input []byte) error {
	fromAddress := from.Address()
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	rawTx := types.NewTransaction(nonce, contract, big.NewInt(0), big.NewInt(0), big.NewInt(0), input)
	signedTx, err := from.SignTransaction(types.HomesteadSigner{}, rawTx)
	if err != nil {
		return err
	}

	return c.SendTransaction(context.Background(), signedTx)
}

func (c *Client) LocalCall(contract common.Address, from *Account, input []byte) ([]byte, error) {
	msg := ethereum.CallMsg{From: from.Address(), To: &contract, Data: input}
	return c.CallContract(context.Background(), msg, nil)
}

func PackFunctionCall(api abi.ABI, method string, args ...interface{}) ([]byte, error) {
	return api.Pack(method, args...)
}

func ABIFromString(api string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(api))
}

func (c *Client) Deploy(from *Account, api, bytecode string, params ...interface{}) (common.Address, error) {
	parsed, err := ABIFromString(api)
	if err != nil {
		return common.Address{}, err
	}
	address, _, _, err := bind.DeployContract(from.Transactor(), parsed, common.FromHex(bytecode), c, params...)
	return address, err
}
