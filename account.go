package ethclient

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	key      *keystore.Key
	password string
}

func NewAccount(keyFile, password string) (*Account, error) {
	key, err := importKey(keyFile, password)
	if err != nil {
		return nil, err
	}

	return &Account{
		key:      key,
		password: password,
	}, nil
}

func importKey(keyFile, password string) (*keystore.Key, error) {
	f, err := os.Open(keyFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	json, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return keystore.DecryptKey(json, password)
}

func (account *Account) Address() common.Address {
	return account.key.Address
}

func (account *Account) SignTransaction(signer types.Signer, tx *types.Transaction) (*types.Transaction, error) {
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), account.key.PrivateKey)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(signer, signature)
}

func (c *Client) AccountNextNonce(address common.Address) (uint64, error) {
	return c.PendingNonceAt(context.Background(), address)
}

func (c *Client) Transfer(from *Account, destAccount common.Address, nonce uint64, value *big.Int) error {
	rawTx := types.NewTransaction(nonce, destAccount, value, nil)
	signedTx, err := from.SignTransaction(types.HomesteadSigner{}, rawTx)
	if err != nil {
		return err
	}

	return c.SendTransaction(context.Background(), signedTx)
}

func (c *Client) GetBalance(destAccount common.Address) (*big.Int, error) {
	return c.BalanceAt(destAccount, nil)
}

func (account *Account) PrivateKey() string {
	return hex.EncodeToString(crypto.FromECDSA(account.key.PrivateKey))
}

func (account *Account) Password() string {
	return account.password
}

func (account *Account) Transactor() *bind.TransactOpts {
	return bind.NewKeyedTransactor(account.key.PrivateKey)
}
