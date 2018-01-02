package ethclient

import (
	"github.com/ethereum/go-ethereum/common"
)

func (ec *Client) Coinbase() (common.Address, error) {
	var address common.Address
	err := ec.c.Call(&address, "eth_coinbase")
	return address, err
}

func (ec *Client) SetCoinbase(account *Account) (bool, error) {
	var address common.Address
	err := ec.c.Call(&address, "personal_importRawKey", account.PrivateKey(), account.Password())
	if err != nil {
		return false, err
	}

	var succeed bool
	err = ec.c.Call(&succeed, "personal_unlockAccount", account.Address(), account.Password())
	if err != nil {
		return succeed, err
	}

	succeed = false
	err = ec.c.Call(&succeed, "miner_setEtherbase", account.Address())
	return succeed, err
}

func (ec *Client) MinerStart(threadCount int) error {
	return ec.c.Call(nil, "miner_start", threadCount)
}

func (ec *Client) MinerStop() (bool, error) {
	var succeed bool
	err := ec.c.Call(&succeed, "miner_stop")
	return succeed, err
}
