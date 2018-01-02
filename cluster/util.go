package cluster

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	extraSeal   = 65
	extraVanity = 32
)

func Ecrecover(header *types.Header) (common.Address, error) {
	signature := header.Extra[len(header.Extra)-extraSeal:]
	pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	return signer, nil
}

func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])
	return hash
}

func headerInfo(b *types.Block) (voteAddr string, isVote bool, signers []string, isEpoch bool) {
	isVote = b.Nonce() == math.MaxUint64
	voteAddr = b.Coinbase().Hex()

	h := b.Header()
	if len(h.Extra) > (extraSeal + extraVanity) {
		signerAddrs := h.Extra[extraVanity:(len(h.Extra) - extraSeal)]
		if len(signerAddrs)%common.AddressLength == 0 {
			signerCount := len(signerAddrs) / common.AddressLength
			isEpoch = true
			for i := 0; i < signerCount; i++ {
				addr := common.BytesToAddress(signerAddrs[i*common.AddressLength : (i+1)*common.AddressLength]).Hex()
				signers = append(signers, addr)
			}
		}
	}
	return
}

type BlockMarshaling struct {
	Signer       string    `json:"signer"`
	Number       uint64    `json:"number"`
	Hash         string    `json:"hash"`
	Parent       string    `json:"parent"`
	Transactions []string  `json:"transactions"`
	Difficulty   uint64    `json:"difficulty"`
	Time         time.Time `json:"time"`
	IsVote       bool      `json:"is_vote"`
	VoteAddress  string    `json:"vote_address"`
	IsEpoch      bool      `json:"is_epoch"`
	Signers      []string  `json:"signers"`
}

func BlockInPOA(block *types.Block, ctrl *Controller) BlockMarshaling {
	h := block.Header()
	signerAddr, _ := Ecrecover(h)

	var txs []string
	for _, tx := range block.Body().Transactions {
		txs = append(txs, tx.Hash().Hex())
	}

	voteAddr, isVote, signers, isEpoch := headerInfo(block)
	return BlockMarshaling{
		Signer:       ctrl.SignerWithAccount(signerAddr),
		Number:       h.Number.Uint64(),
		Hash:         h.Hash().Hex(),
		Parent:       h.ParentHash.Hex(),
		Transactions: txs,
		Difficulty:   h.Difficulty.Uint64(),
		Time:         time.Unix(h.Time.Int64(), 0),
		IsVote:       isVote,
		VoteAddress:  voteAddr,
		IsEpoch:      isEpoch,
		Signers:      signers,
	}
}

type TransactionMarshaling struct {
	Hash         string `json:"hash"`
	AccountNonce uint64 `json:"account_nonce"`
	From         string `json:"from"`
	To           string `json:"to"`
	Value        uint64 `json:"value"`
	Data         []byte `json:"data"`
}

func TransactionInPOA(tx *types.Transaction) TransactionMarshaling {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		panic("invalid tx:" + err.Error())
	}

	toAddress := msg.To()
	var to string
	if toAddress != nil {
		to = toAddress.Hex()
	}

	return TransactionMarshaling{
		Hash:         tx.Hash().Hex(),
		AccountNonce: msg.Nonce(),
		From:         msg.From().Hex(),
		To:           to,
		Value:        msg.Value().Uint64(),
		Data:         tx.Data(),
	}
}

type ReceiptMarshaling struct {
	Succeed bool            `json:"succeed"`
	GasUsed uint64          `json:"gas_used"`
	Logs    []LogMarshaling `json:"logs"`
}

func ReceiptInPOA(r *types.Receipt) ReceiptMarshaling {
	var logs []LogMarshaling
	for _, log := range r.Logs {
		logs = append(logs, LogInPOA(log))
	}
	return ReceiptMarshaling{
		Succeed: r.Status == types.ReceiptStatusSuccessful,
		GasUsed: r.CumulativeGasUsed.Uint64(),
		Logs:    logs,
	}
}

type LogMarshaling struct {
	Topics string `json:"topics"`
	Data   string `json:"data"`
}

func LogInPOA(l *types.Log) LogMarshaling {
	return LogMarshaling{
		Topics: fmt.Sprintf("%x", l.Topics),
		Data:   fmt.Sprintf("%x", l.Data),
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
