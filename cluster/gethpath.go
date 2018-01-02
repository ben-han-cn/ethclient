package cluster

import (
	"fmt"
	"path/filepath"
)

type GethPath struct {
	datapath string
	gethpath string
}

func NewGethPath(conf *Config) (*GethPath, error) {
	datapath, err := filepath.Abs(conf.NodeDataPath)
	if err != nil {
		return nil, err
	}

	gethpath, err := filepath.Abs(conf.GethPath)
	if err != nil {
		return nil, err
	}

	return &GethPath{
		datapath: datapath,
		gethpath: gethpath,
	}, nil
}

func (gp *GethPath) NodeDataPath(n *Node) string {
	return filepath.Join(gp.datapath, n.Name())
}

func (gp *GethPath) NodeLogPath(n *Node) string {
	return filepath.Join(gp.datapath, fmt.Sprintf("%s.log", n.Name()))
}

func (gp *GethPath) GethPath() string {
	return gp.gethpath
}

func (gp *GethPath) GenesisPath() string {
	return filepath.Join(gp.datapath, "poa_for_fun.json")
}
