package cluster

import (
	"context"
	"math/big"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"goldenteam/ethclient"
)

const DefaultPasswd = "ge1888&x399^%$xx__999xx13K"

type Controller struct {
	nodeManager    *NodeManager
	signerAccounts map[string]common.Address
	signers        *nodeSlice
	syncers        *nodeSlice
	keyGenerator   *ethclient.KeyGenerator
}

func NewController(conf *Config) (*Controller, error) {
	nodeManager, err := NewNodeManager(conf)
	if err != nil {
		return nil, err
	}

	keyFilePath, err := filepath.Abs(conf.KeyStorePath)
	if err != nil {
		return nil, err
	}

	return &Controller{
		nodeManager:    nodeManager,
		signerAccounts: make(map[string]common.Address),
		signers:        NewNodeSlice(Signer, conf.SignerCount),
		syncers:        NewNodeSlice(Syncer, conf.SyncerCount),
		keyGenerator:   ethclient.NewKeyGenerator(keyFilePath),
	}, nil
}

func (c *Controller) Signers() *nodeSlice { return c.signers }
func (c *Controller) Syncers() *nodeSlice { return c.syncers }
func (c *Controller) SignerWithAccount(target common.Address) string {
	for name, address := range c.signerAccounts {
		if address == target {
			return name
		}
	}
	return ""
}

func (c *Controller) ClientForNode(n *Node) (*ethclient.Client, error) {
	if c.signers.Include(n) == false && c.syncers.Include(n) == false {
		return nil, ErrUnknownNode
	}
	return c.nodeManager.Client(n), nil
}

func (c *Controller) StartAllNodes() error {
	signerNodes := c.signers.Nodes()
	if err := c.generateGensisBlock(); err != nil {
		return err
	}

	syncerNodes := c.syncers.Nodes()
	if err := c.nodeManager.StartNodes(append(signerNodes, syncerNodes...)); err != nil {
		return err
	}

	for _, signer := range signerNodes {
		if err := c.startSign(signer); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) StopNode(n *Node) error {
	if n.Role() == Signer {
		if err := c.signers.RemoveNode(n); err != nil {
			return err
		}

		for _, otherSigner := range c.signers.Nodes() {
			c.nodeManager.Client(otherSigner).Propose(c.signerAccounts[n.Name()], false)
		}

		return nil
	} else {
		if err := c.syncers.RemoveNode(n); err != nil {
			return err
		}
		c.syncers.RemoveNode(n)
		return c.nodeManager.StopNode(n)
	}
}

func (c *Controller) RestartNode(n *Node) error {
	if n.Role() == Signer {
		otherSigners := c.signers.Nodes()
		if err := c.signers.RestoreNode(n); err != nil {
			return err
		}

		for _, otherSigner := range otherSigners {
			c.nodeManager.Client(otherSigner).Propose(c.signerAccounts[n.Name()], true)
		}

		return nil
	} else {
		if err := c.syncers.RestoreNode(n); err != nil {
			return err
		}
		return c.nodeManager.StartNode(n)
	}
}

func (c *Controller) AddSyncer() (*Node, error) {
	syncer, err := c.syncers.AddNode()
	if err != nil {
		return nil, err
	}

	if err := c.nodeManager.StartNode(syncer); err == nil {
		return syncer, nil
	} else {
		return nil, err
	}
}

func (c *Controller) AddSigner() (*Node, error) {
	runningSigners := c.signers.Nodes()

	signer, err := c.signers.AddNode()
	if err != nil {
		return nil, err
	}

	accounts, err := c.createAccount(1)
	if err != nil {
		return nil, err
	}

	c.signerAccounts[signer.Name()] = accounts[0]
	if err := c.nodeManager.StartNode(signer); err != nil {
		return nil, err
	}

	if err := c.startSign(signer); err != nil {
		return nil, err
	}

	for _, otherSigner := range runningSigners {
		c.nodeManager.Client(otherSigner).Propose(accounts[0], true)
	}
	return signer, nil
}

func (c *Controller) generateGensisBlock() error {
	blockDuration := time.Duration(5) * time.Second

	accounts, err := c.createAccount(c.signers.Size())
	if err != nil {
		return err
	}

	testAccounts, err := c.createAccount(10)
	funds := make(map[common.Address]*big.Int)
	for _, account := range append(accounts, testAccounts...) {
		funds[account] = big.NewInt(1000000000000000000)
	}

	genesis := ethclient.MakeGenesis(blockDuration, accounts, funds)
	if err := ethclient.CreateGensisFile(genesis, c.nodeManager.GenesisPath()); err != nil {
		return err
	}

	for i, n := range c.signers.Nodes() {
		c.signerAccounts[n.Name()] = accounts[i]
	}
	return nil
}

func (c *Controller) createAccount(n int) ([]common.Address, error) {
	var addresses []common.Address
	for i := 0; i < n; i++ {
		addr, err := c.keyGenerator.GenerateKey(DefaultPasswd)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

func (c *Controller) startSign(signer *Node) error {
	client := c.nodeManager.Client(signer)
	address := c.signerAccounts[signer.Name()]
	account, err := c.keyGenerator.GetAccount(address, DefaultPasswd)
	if err != nil {
		return err
	}
	client.SetCoinbase(account)
	return client.MinerStart(1)
}

func (c *Controller) Accounts() []common.Address {
	return c.keyGenerator.ListAddress()
}

func (c *Controller) TransferMoney(from, to common.Address, value int64, count int) error {
	runningSigners := c.Signers().Nodes()
	runningSyncers := c.Syncers().Nodes()
	syncerName := runningSyncers[rand.Intn(len(runningSyncers))]
	signerName := runningSigners[rand.Intn(len(runningSigners))]

	fromAccount, err := c.keyGenerator.GetAccount(from, DefaultPasswd)
	if err != nil {
		return err
	}
	nonce, err := c.nodeManager.Client(signerName).PendingNonceAt(context.Background(), fromAccount.Address())
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = c.nodeManager.Client(syncerName).Transfer(fromAccount, to, nonce, big.NewInt(value))
		if err != nil {
			return err
		}
		nonce += 1
	}
	return nil
}
