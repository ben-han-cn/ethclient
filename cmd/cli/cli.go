package main

import (
	"flag"

	"github.com/c-bata/go-prompt"
	"goldenteam/ethclient/cluster"
)

var (
	nodeDataPath string
	keyStore     string
	signerCount  int
	syncerCount  int
	gethPath     string
)

func init() {
	flag.StringVar(&nodeDataPath, "n", "", "top folder for all nodes")
	flag.StringVar(&keyStore, "k", "", "folder to store all the key files")
	flag.StringVar(&gethPath, "p", "", "geth cmd path")
	flag.IntVar(&signerCount, "s", 2, "signer node count")
	flag.IntVar(&syncerCount, "y", 3, "syncer node count")
}

func main() {
	flag.Parse()
	ctrl, err := cluster.NewController(&cluster.Config{
		NodeDataPath: nodeDataPath,
		KeyStorePath: keyStore,
		SignerCount:  signerCount,
		SyncerCount:  syncerCount,
		GethPath:     gethPath,
	})
	if err != nil {
		panic("create ctrl failed:" + err.Error())
	}

	if err := ctrl.StartAllNodes(); err != nil {
		panic("start nodes failed:" + err.Error())
	}

	p := prompt.New(
		NewExecutor(ctrl),
		NewCompleter(ctrl),
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("ethctrl-prompt"),
	)
	p.Run()
}
