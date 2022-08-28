package block

import (
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/common"
	"jokecoin-go/core/utils"
)

// Transaction represents a transaction
// TxType: 1 - normal TX, 2 - creating block TX
type Transaction struct {
	TxType          int64            `json:"tx_type"`
	SenderPublicKey common.PublicKey `json:"sender_public_key"`
	SenderSignature common.Signature `json:"sender_signature"`
	Receiver        common.Address   `json:"receiver"`
	Value           int64            `json:"value"`
	GasLimit        int64            `json:"gas_limit"`
	Fee             int64            `json:"fee"`
	Nonce           int64            `json:"nonce"`
	Data            []byte           `json:"data"`
}

func (tx Transaction) Encode() []byte {
	b, err := utils.Serialize(tx)
	utils.PanicIfErr(err)
	return b
}

func (tx Transaction) ComputeHash() common.Hash {
	return sha3.Sum256(tx.Encode())
}
