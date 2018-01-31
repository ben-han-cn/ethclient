package cluster

import (
	"os"
	"strconv"
	"sync"
)

type gethRunner struct {
	gethpath    *GethPath
	nextp2pport int
	nextrpcport int
	portMu      sync.Mutex
}

func NewGethRunner(gethpath *GethPath) *gethRunner {
	return &gethRunner{
		gethpath:    gethpath,
		nextp2pport: 8800,
		nextrpcport: 9900,
	}
}

func (r *gethRunner) StartGeth(n *Node) (*os.Process, error) {
	nodedatapath := r.gethpath.NodeDataPath(n)
	if PathExists(nodedatapath) == false {
		if err := os.MkdirAll(nodedatapath, 0700); err != nil {
			return nil, err
		}

		if err := r.initGenesis(nodedatapath); err != nil {
			return nil, err
		}
	}

	p2pport, rpcport := r.allocatePort()
	return startProcess(r.gethpath.GethPath(), r.gethpath.NodeLogPath(n),
		"--datadir", nodedatapath,
		"--networkid", "77877",
		"--nodiscover",
		"--rpc",
		"--rpcport", strconv.Itoa(rpcport),
		"--port", strconv.Itoa(p2pport))
}

func (r *gethRunner) initGenesis(nodedatapath string) error {
	p, err := startProcess(r.gethpath.GethPath(), "",
		"--datadir", nodedatapath,
		"init", r.gethpath.GenesisPath())

	if err != nil {
		return err
	}

	p.Wait()
	return nil
}

func (r *gethRunner) allocatePort() (p2pport int, rpcport int) {
	r.portMu.Lock()
	defer r.portMu.Unlock()

	p2pport = r.nextp2pport
	rpcport = r.nextrpcport
	r.nextp2pport += 1
	r.nextrpcport += 1
	return
}

func startProcess(cmd string, outputFile string, cmd_args ...string) (*os.Process, error) {
	var stdout, stderr *os.File
	if outputFile != "" {
		out, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		stdout = out
		stderr = out
	}

	attr := &os.ProcAttr{
		Files: []*os.File{nil, stdout, stderr},
	}

	args := append([]string{cmd}, cmd_args...)
	return os.StartProcess(cmd, args, attr)
}
