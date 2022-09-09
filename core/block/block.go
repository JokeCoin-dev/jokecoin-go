package block

import (
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/common"
	"jokecoin-go/core/database"
	"jokecoin-go/core/errors"
	"jokecoin-go/core/utils"
	"log"
)

const MaxExtraDataSize = 64
const MAX_TXS_PER_BLOCK = 1000

type BlockHeader struct {
	//Hash            common.Hash
	ParentHash       common.Hash `json:"parent_hash"`
	TransactionHash  common.Hash `json:"transaction_hash"`
	StateHash        common.Hash `json:"state_hash"`
	Time             int64       `json:"time"`
	ExtraData        []byte      `json:"extra_data"`
	Difficulty       []byte      `json:"difficulty"` // convert to big.Int to use it
	Height           int64       `json:"height"`
	LastBlockTime    int64       `json:"last_block_time"`
	LastKeyBlockTime int64       `json:"last_key_block_time"`
}

type Block struct {
	Header BlockHeader   `json:"header"`
	TXs    []Transaction `json:"txs"`
}

func (h BlockHeader) ComputeHash() common.Hash {
	utils.Assert(h.Time >= 0)
	hash, err := utils.Serialize(h)
	utils.PanicIfErr(err)
	return sha3.Sum256(hash)
}

func LoadBlock(hash common.Hash) (*Block, error) {
	b := &Block{}
	data, err := database.GetDB().Get(hash[:])
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	err = utils.Deserialize(data, b)
	if err != nil {
		return nil, errors.Wrap(err, errors.DataConsistencyError)
	}
	return b, err
}

func (b Block) WriteDB() error {
	if b.TXs == nil || (len(b.TXs) == 0 && b.Header.TransactionHash != common.EmptyHash) {
		log.Fatalln("Try to encode an uncompleted block")
	}
	hash := b.Header.ComputeHash()
	data, err := utils.Serialize(b)
	if err != nil {
		return err
	}
	return database.GetDB().Put(hash[:], data)
}

func CheckGenesisBlock(block *Block) bool {
	if block.Header.ParentHash != common.EmptyHash {
		return false
	}
	if block.Header.TransactionHash != common.EmptyHash {
		return false
	}
	if block.Header.StateHash != common.EmptyHash {
		return false
	}
	if block.Header.Time != 0 {
		return false
	}
	if block.Header.Height != 0 {
		return false
	}
	if block.Header.LastBlockTime != 0 || block.Header.LastKeyBlockTime != 0 {
		return false
	}
	if len(block.Header.Difficulty) != common.HashLen {
		return false
	}
	if len(block.TXs) != 0 {
		return false
	}
	if len(block.Header.ExtraData) > MaxExtraDataSize {
		return false
	}
	return true
}
