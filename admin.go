package ethclient

import (
	"github.com/ethereum/go-ethereum/p2p"
)

func (ec *Client) AddPeer(url string) (bool, error) {
	var succeed bool
	err := ec.c.Call(&succeed, "admin_addPeer", url)
	if err != nil {
		return false, err
	} else {
		return succeed, nil
	}
}

func (ec *Client) Peers() ([]*p2p.PeerInfo, error) {
	var peers []*p2p.PeerInfo
	err := ec.c.Call(&peers, "admin_peers")
	if err != nil {
		return nil, err
	} else {
		return peers, nil
	}
}

func (ec *Client) NodeInfo() (*p2p.NodeInfo, error) {
	var nodeInfo p2p.NodeInfo
	err := ec.c.Call(&nodeInfo, "admin_nodeInfo")
	if err != nil {
		return nil, err
	} else {
		return &nodeInfo, nil
	}
}
