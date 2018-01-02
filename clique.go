package ethclient

import (
	"github.com/ethereum/go-ethereum/common"
)

func (ec *Client) Propose(address common.Address, auth bool) error {
	return ec.c.Call(nil, "clique_propose", address, auth)
}

func (ec *Client) Proposes() map[common.Address]bool {
	var proposes map[common.Address]bool
	ec.c.Call(proposes, "clique_proposes")
	return proposes
}
