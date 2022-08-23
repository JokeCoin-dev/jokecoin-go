package block

import (
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/common"
	"jokecoin-go/core/utils"
	"math/big"
)

type BlockHeader struct {
	//Hash            common.Hash
	ParentHash      common.Hash
	TransactionHash common.Hash
	StateHash       common.Hash
	Difficulty      *big.Int
	Time            int64
	ExtraData       []byte
}

type Block struct {
	Header BlockHeader
	TXs    []Transaction
}

func (h *BlockHeader) ComputeHash() common.Hash {
	utils.Assert(h.Time >= 0)
	hash, err := utils.Serialize(h)
	utils.PanicIfErr(err)
	return sha3.Sum256(hash)
}

func (b *Block) EncodeBlockBody() []byte {
	buf := make([]byte, common.AddressLen+8)
	return buf
}
