package main

import (
	"flag"
	"log"

	"github.com/c-bata/go-prompt"
	"goldenteam/ethclient/cluster"
)

var (
	nodeDataPath    string
	keyStore        string
	contractAddress string
)

func init() {
	flag.StringVar(&nodeDataPath, "n", "", "top folder for all nodes")
	flag.StringVar(&keyStore, "k", "", "folder to store all the key files")
	flag.StringVar(&contractAddress, "c", "", "contract address which has been deployed")
}

func main() {
	flag.Parse()

	client, err := newContractClient(&cluster.Config{
		NodeDataPath: nodeDataPath,
		KeyStorePath: keyStore,
	}, contractAddress)

	if err != nil {
		log.Fatalf("connect to fullnode failed:%s", err.Error())
	}

	p := prompt.New(
		NewExecutor(client),
		NewCompleter(client),
		prompt.OptionPrefix(client.CurrentNodeName()+">>> "),
		prompt.OptionTitle("contract-tester"),
	)
	p.Run()
}
