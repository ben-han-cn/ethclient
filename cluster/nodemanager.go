package cluster

import (
	"log"
	"os"
	"sync"

	"goldenteam/ethclient"
)

type NodeManager struct {
	clients   map[string]*ethclient.Client
	processes map[string]*os.Process
	mu        sync.Mutex
	runner    *gethRunner
	gethpath  *GethPath
}

func NewNodeManager(conf *Config) (*NodeManager, error) {
	gethpath, err := NewGethPath(conf)
	if err != nil {
		return nil, err
	}

	return &NodeManager{
		clients:   make(map[string]*ethclient.Client),
		processes: make(map[string]*os.Process),
		runner:    NewGethRunner(gethpath),
		gethpath:  gethpath,
	}, nil
}

func (nm *NodeManager) Client(n *Node) *ethclient.Client {
	client, ok := nm.clients[n.Name()]
	if ok == false {
		panic("get unknown node" + n.Name())
	}
	return client
}

func (nm *NodeManager) StartNodes(nodes []*Node) error {
	errChan := make(chan error, len(nodes))
	for _, node := range nodes {
		go func(n *Node) {
			errChan <- nm.StartNode(n)
		}(node)
	}

	for i := 0; i < len(nodes); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func (nm *NodeManager) StartNode(n *Node) error {
	if nm.HasNode(n) {
		return nil
	}

	log.Printf("start to launch node %s\n", n.Name())
	p, err := nm.runner.StartGeth(n)
	if err != nil {
		return err
	}

	client := n.Client(nm.gethpath)

	nm.mu.Lock()
	nm.processes[n.Name()] = p
	for _, other := range nm.clients {
		other.AddPeer(client.Enode())
	}
	nm.clients[n.Name()] = client
	nm.mu.Unlock()
	log.Printf("finish launch node %s\n", n.Name())
	return nil
}

func (nm *NodeManager) HasNode(n *Node) bool {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	_, ok := nm.clients[n.Name()]
	return ok
}

const connectRetry = 10

func (nm *NodeManager) StopNode(n *Node) error {
	nodeName := n.Name()
	client := nm.clients[nodeName]
	if client == nil {
		return nil
	}

	process := nm.processes[nodeName]
	process.Kill()
	process.Wait()

	client.Close()
	delete(nm.clients, nodeName)
	delete(nm.processes, nodeName)
	return nil
}

func (nm *NodeManager) GenesisPath() string {
	return nm.gethpath.GenesisPath()
}
