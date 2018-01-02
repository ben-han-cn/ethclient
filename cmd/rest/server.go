package main

import (
	"errors"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"cement/httprouter"
	"github.com/ethereum/go-ethereum/common"
	"goldenteam/ethclient"
	"goldenteam/ethclient/cluster"
)

var (
	ErrUnknownRole = errors.New("unknown role")
)

const DefaultPassword string = "120987_ABC$^^^*acbcedd"

type Server struct {
	ctrl *cluster.Controller
}

func NewServer(ctrl *cluster.Controller) *Server {
	return &Server{
		ctrl: ctrl,
	}
}

func (s *Server) Run(addr string) {
	rand.Seed(time.Now().UnixNano())

	router := httprouter.New()
	router.GET("/api/v1/nodes", s.getNodes)
	router.POST("/api/v1/nodes", s.addNode)

	router.GET("/api/v1/nodes/:node/blocks/:number", s.getBlock)
	router.GET("/api/v1/nodes/:node/transactions/:hash", s.getTransaction)
	router.POST("/api/v1/transactions", s.transferMoney)
	router.GET("/api/v1/accounts", s.getAccounts)

	if err := http.ListenAndServe(addr, router); err != nil {
		panic("start server failed:" + err.Error())
	}
}

func (s *Server) getNodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var signers []string
	for _, n := range s.ctrl.Signers().Nodes() {
		signers = append(signers, n.Name())
	}

	var syncers []string
	for _, n := range s.ctrl.Syncers().Nodes() {
		syncers = append(syncers, n.Name())
	}

	EncodeResult(w, SucceedWithResult(&struct {
		Signers []string `json:"signers"`
		Syncers []string `json:"syncers"`
	}{
		signers,
		syncers,
	}))
}

func (s *Server) addNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	param := struct {
		Role string `json:"role"`
	}{}
	err := DecodeRequestBody(r, &param)
	if err != nil {
		EncodeResult(w, Failed(ErrInvalidParameter))
		return
	}

	var newNode *cluster.Node
	switch cluster.Role(param.Role) {
	case cluster.Signer:
		newNode, err = s.ctrl.AddSigner()
	case cluster.Syncer:
		newNode, err = s.ctrl.AddSyncer()
	default:
		err = ErrUnknownRole
	}

	if err != nil {
		EncodeResult(w, Failed(ErrInvalidParameter))
		return
	} else {
		EncodeResult(w, SucceedWithResult(
			struct {
				Id   string `json:"id"`
				Role string `json:"role"`
			}{
				Id:   newNode.Name(),
				Role: param.Role,
			}))
	}
}

func (s *Server) connectToNode(name string) (*ethclient.Client, error) {
	node, err := cluster.NodeFromString(name)
	if err != nil {
		return nil, err
	}

	return s.ctrl.ClientForNode(node)
}

func (s *Server) getBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	node, err := s.connectToNode(ps.ByName("node"))
	if err != nil {
		EncodeResult(w, Failed(ErrUnknownNode))
		return
	}

	number, _ := strconv.Atoi(ps.ByName("number"))
	var blockNumber *big.Int
	if number >= 0 {
		blockNumber = big.NewInt(int64(number))
	}

	block, err := node.BlockByNumber(blockNumber)
	if err != nil {
		EncodeResult(w, Failed(ErrGetBlockFailed))
	} else {
		EncodeResult(w, SucceedWithResult(cluster.BlockInPOA(block, s.ctrl)))
	}
}

func (s *Server) getTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	node, err := s.connectToNode(ps.ByName("node"))
	if err != nil {
		EncodeResult(w, Failed(ErrUnknownNode))
		return
	}

	tx, _, err := node.TransactionByHash(common.HexToHash(ps.ByName("hash")))
	if err != nil {
		EncodeResult(w, Failed(ErrGetTransactionFailed))
	} else {
		EncodeResult(w, SucceedWithResult(cluster.TransactionInPOA(tx)))
	}
}

func (s *Server) transferMoney(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	param := struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Value int64  `json:"value"`
		Count int    `json:"count"`
	}{}
	err := DecodeRequestBody(r, &param)
	if err != nil {
		EncodeResult(w, Failed(ErrInvalidParameter))
		return
	}

	fromAccount := common.HexToAddress(param.From)
	toAccount := common.HexToAddress(param.To)
	s.ctrl.TransferMoney(fromAccount, toAccount, param.Value, param.Count)
	EncodeResult(w, Succeed())
}

func (s *Server) getAccounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	accounts := s.ctrl.Accounts()
	addresses := make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		addresses[i] = accounts[i].Hex()
	}
	EncodeResult(w, SucceedWithResult(struct {
		Accounts []string `json:"accounts"`
	}{addresses}))
}
