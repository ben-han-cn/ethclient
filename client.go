package ethclient

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	c *rpc.Client
}

func Dial(rawurl string) (*Client, error) {
	c, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}

func (ec *Client) Close() {
	ec.c.Close()
}

func NewIPCClient(dataPath string) (*Client, error) {
	return Dial(ipcPath(dataPath))
}

func ipcPath(dataPath string) string {
	return filepath.Join(dataPath, "geth.ipc")
}

func (c *Client) Enode() string {
	info, err := c.NodeInfo()
	if err != nil {
		log.Fatalf("get nodeinfo failed:%s", err.Error())
	}
	return strings.Replace(info.Enode, "[::]", "127.0.0.1", 1)
}
