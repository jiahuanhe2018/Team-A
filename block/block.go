package block

import (
	TX "../transaction"
	"../pow"
	"../utils"
	"../account"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index        int                        `json:"index"`
	Timestamp    string                     `json:"timestamp"`
	Result       int                        `json:"result"`
	Hash         string                     `json:"hash"`
	PrevHash     string                     `json:"prevhash"`
	Proof        uint64                     `json:"proof"`
	Transactions []TX.Transaction           `json:"transactions"`
	Accounts     map[string]account.Account `json:"accounts"`
	Difficulty   int                        `json:"difficulty"`
	Nonce        string                     `json:"nonce"`
	Validator    string                     `json:"validator"`
}

func GenerateBlock(oldBlock *Block, Result int, address string) *Block {
	newBlock := &Block{
		Index:     oldBlock.Index + 1,
		Timestamp: time.Now().String(),
		Result:    Result,
		Hash:      "",
		PrevHash:  oldBlock.Hash,
	}

	switch utils.Consensus {
	case utils.ConsensusPow:
		newBlock.Difficulty = utils.PowDifficulty
		pow.MineHash(newBlock)
	case utils.ConsensusPoS:
		newBlock.Validator = address
		newBlock.Hash = newBlock.CalculateHash()
	}
	return newBlock
}

func (b *Block) IsBlockValid(oldBlock *Block) bool {
	if oldBlock.Index+1 != b.Index {
		return false
	}
	if oldBlock.Hash != b.PrevHash {
		return false
	}
	if b.CalculateHash() != b.Hash {
		return false
	}
	return true
}

func (b *Block) CalculateHash() string {
	record := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.Result) + b.PrevHash + b.Nonce
	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}
