package cluster

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"goldenteam/ethclient"
)

var (
	ErrNodeNameIsNotValid = errors.New("node name isn't valid which should only be signerN or syncerN and N start from 1")
	ErrUnknownNode        = errors.New("node is unknown")
	ErrRoleUnMatch        = errors.New("node isn't belongs to current slice it has different role")
	ErrNodeAlreadyExists  = errors.New("node already exists")
	ErrTooManyNode        = errors.New("too many nodes")
)

type Role string

var (
	Signer Role = "signer"
	Syncer Role = "syncer"
)

type Node struct {
	role  Role
	index int
}

func (n *Node) Index() int           { return n.index }
func (n *Node) Name() string         { return fmt.Sprintf("%s%d", string(n.role), n.index+1) }
func (n *Node) Role() Role           { return n.role }
func (n *Node) Equals(on *Node) bool { return n.role == on.role && n.index == on.index }

func NodeFromString(name string) (*Node, error) {
	var index int
	var role Role

	if strings.HasPrefix(name, string(Signer)) {
		role = Signer
	} else if strings.HasPrefix(name, string(Syncer)) {
		role = Syncer
	} else {
		return nil, ErrNodeNameIsNotValid
	}

	index, err := strconv.Atoi(strings.TrimPrefix(name, string(role)))
	if err != nil {
		return nil, err
	}

	if index < 1 {
		return nil, ErrNodeNameIsNotValid
	}

	return &Node{
		role:  role,
		index: index - 1,
	}, nil
}

func (n *Node) Client(path *GethPath) *ethclient.Client {
	retryCount := 0
	for {
		node, err := ethclient.NewIPCClient(path.NodeDataPath(n))
		if err == nil {
			return node
		}
		if retryCount += 1; retryCount == connectRetry {
			log.Fatalf("get node to %s failed %s", n.Name(), err.Error())
		}
		<-time.After(time.Second)
	}
	return nil
}

const maxNodeCount = 100

type nodeSlice struct {
	role   Role
	exists [maxNodeCount]bool
	count  int
}

func NewNodeSlice(role Role, count int) *nodeSlice {
	if count > maxNodeCount {
		panic("too many nodes")
	}

	ns := &nodeSlice{role: role, count: count}
	for i := 0; i < count; i++ {
		ns.exists[i] = true
	}
	return ns
}

func (ns *nodeSlice) Size() int {
	return ns.count
}

func (ns *nodeSlice) Include(n *Node) bool {
	return n.role == ns.role && n.Index() < ns.count && ns.exists[n.Index()]
}

func (ns *nodeSlice) RemoveNode(n *Node) error {
	if n.Role() != ns.role {
		return ErrRoleUnMatch
	} else if n.Index() >= ns.count || ns.exists[n.Index()] == false {
		return ErrUnknownNode
	}

	ns.exists[n.Index()] = false
	return nil
}

func (ns *nodeSlice) RestoreNode(n *Node) error {
	if n.Role() != ns.role {
		return ErrRoleUnMatch
	} else if n.Index() >= ns.count {
		return ErrUnknownNode
	} else if ns.exists[n.Index()] == true {
		return ErrNodeAlreadyExists
	}

	ns.exists[n.Index()] = true
	return nil
}

func (ns *nodeSlice) AddNode() (*Node, error) {
	if ns.count+1 >= maxNodeCount {
		return nil, ErrTooManyNode
	}
	ns.exists[ns.count] = true
	ns.count += 1
	return &Node{role: ns.role, index: ns.count - 1}, nil
}

func (ns *nodeSlice) MissingNodes() (nodes []*Node) {
	for i := 0; i < ns.count; i++ {
		if ns.exists[i] == false {
			nodes = append(nodes, &Node{role: ns.role, index: i})
		}
	}
	return
}

func (ns *nodeSlice) Nodes() (nodes []*Node) {
	for i := 0; i < ns.count; i++ {
		if ns.exists[i] {
			nodes = append(nodes, &Node{role: ns.role, index: i})
		}
	}
	return
}
