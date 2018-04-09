package ethclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
)

func MakeGenesis(blockInterval time.Duration, signers []common.Address, prefund map[common.Address]*big.Int) *core.Genesis {
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: big.NewInt(1048576),
		Alloc:      make(core.GenesisAlloc),
		Config:     &params.ChainConfig{},
	}
	genesis.Difficulty = big.NewInt(1)
	genesis.Config.Clique = &params.CliqueConfig{
		Period: uint64(blockInterval.Seconds()),
		Epoch:  30,
	}

	for i := 0; i < len(signers); i++ {
		for j := i + 1; j < len(signers); j++ {
			if bytes.Compare(signers[i][:], signers[j][:]) > 0 {
				signers[i], signers[j] = signers[j], signers[i]
			}
		}
	}
	genesis.ExtraData = make([]byte, 32+len(signers)*common.AddressLength+65)
	for i, signer := range signers {
		copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
	}

	for address, fund := range prefund {
		genesis.Alloc[address] = core.GenesisAccount{
			Balance: fund,
		}
	}
	// Add a batch of precompile balances to avoid them getting deleted
	for i := int64(0); i < 256; i++ {
		genesis.Alloc[common.BigToAddress(big.NewInt(i))] = core.GenesisAccount{Balance: big.NewInt(1)}
	}

	genesis.Config.ChainId = new(big.Int).SetUint64(uint64(rand.Intn(65536)))
	return genesis
}

func CreateGensisFile(genesis *core.Genesis, fileName string) error {
	out, _ := json.MarshalIndent(genesis, "", "  ")
	return ioutil.WriteFile(fileName, out, 0644)
}
