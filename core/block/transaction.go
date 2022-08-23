package block

import (
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/common"
	"jokecoin-go/core/utils"
)

type Transaction struct {
	TxType          int64
	SenderPublicKey common.PublicKey
	SenderSignature common.Signature
	Receiver        common.Address
	Value           int64
	GasLimit        int64
	Fee             int64
	Nonce           int64
	Data            []byte
}

func (tx Transaction) Encode() []byte {
	b, err := utils.Serialize(tx)
	utils.PanicIfErr(err)
	return b
}

func (tx Transaction) ComputeHash() common.Hash {
	return sha3.Sum256(tx.Encode())
}
