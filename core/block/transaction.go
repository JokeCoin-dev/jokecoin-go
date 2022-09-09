package block

import (
	"crypto/ed25519"
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/common"
	"jokecoin-go/core/errors"
	"jokecoin-go/core/mpt"
	"jokecoin-go/core/utils"
)

const MAX_TX_FEE int64 = 1 << 50
const MAX_TX_VALUE int64 = 1 << 50
const COINBASE_VALUE int64 = 100000000

// Transaction represents a transaction
// TxType: 1 - normal TX, 2 - creating block TX
type Transaction struct {
	TxType          int64             `json:"tx_type"`
	SenderPublicKey ed25519.PublicKey `json:"sender_public_key"`
	SenderSignature common.Signature  `json:"sender_signature"`
	Receiver        common.Address    `json:"receiver"`
	Value           int64             `json:"value"`
	GasLimit        int64             `json:"gas_limit"`
	Fee             int64             `json:"fee"`
	Nonce           int64             `json:"nonce"`
	Data            []byte            `json:"data"`
}

func (tx Transaction) Encode() []byte {
	b, err := utils.Serialize(tx)
	utils.PanicIfErr(err)
	return b
}

func (tx Transaction) ComputeHash() common.Hash {
	return sha3.Sum256(tx.Encode())
}

func (tx *Transaction) Sign(privateKey ed25519.PrivateKey) {
	for i := 0; i < common.SignatureLen; i++ {
		tx.SenderSignature[i] = 0
	}
	t := ed25519.Sign(privateKey, tx.Encode())
	copy(tx.SenderSignature[:], t)
}

func CheckSignature(tx Transaction) bool {
	signature := make([]byte, common.SignatureLen)
	copy(signature[:], tx.SenderSignature[:])
	for i := 0; i < common.SignatureLen; i++ {
		tx.SenderSignature[i] = 0
	}
	res := ed25519.Verify(tx.SenderPublicKey, tx.Encode(), signature[:])
	return res
}

// ExecuteTransaction executes a transaction then return a new state
func ExecuteTransaction(tx Transaction, state mpt.MerklePatriciaTrie) (mpt.MerklePatriciaTrie, error) {
	if tx.TxType == 1 {
		if !CheckSignature(tx) {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Signature invalid")
		}
		senderAddr := common.PublicKeyToAddress(tx.SenderPublicKey)
		b, err := state.Get(senderAddr[:])
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.Wrap(err, errors.DatabaseError)
		}
		if b == nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Sender not found")
		}
		sender := common.AccountData{}
		err = utils.Deserialize(b, &sender)
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.Wrap(err, errors.DataConsistencyError)
		}
		b, err = state.Get(tx.Receiver[:])
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.Wrap(err, errors.DatabaseError)
		}
		receiver := common.AccountData{}
		if b == nil {
			receiver = common.AccountData{Balance: 0, Nonce: 0, CodeHash: common.EmptyHash, StorageHash: common.EmptyHash}
		} else {
			err = utils.Deserialize(b, &receiver)
			if err != nil {
				return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.Wrap(err, errors.DataConsistencyError)
			}
		}
		if tx.Nonce != sender.Nonce+1 {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Nonce not match")
		}
		if tx.Value > MAX_TX_VALUE {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Value too big")
		}
		if tx.Fee > MAX_TX_FEE {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Fee too big")
		}
		if tx.Value+tx.Fee > sender.Balance {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Balance not enough")
		}
		sender.Balance -= tx.Value + tx.Fee
		sender.Nonce++
		receiver.Balance += tx.Value
		state, err = state.Put(senderAddr[:], utils.MustSerialize(sender))
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, err
		}
		state, err = state.Put(tx.Receiver[:], utils.MustSerialize(receiver))
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, err
		}
	} else if tx.TxType == 2 {
		b, err := state.Get(tx.Receiver[:])
		receiver := common.AccountData{}
		if b == nil {
			receiver = common.AccountData{Balance: 0, Nonce: 0, CodeHash: common.EmptyHash, StorageHash: common.EmptyHash}
		} else {
			err = utils.Deserialize(b, &receiver)
			if err != nil {
				return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.Wrap(err, errors.DataConsistencyError)
			}
		}
		if tx.Value != COINBASE_VALUE {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Coinbase value not match")
		}
		receiver.Balance += tx.Value
		state, err = state.Put(tx.Receiver[:], utils.MustSerialize(receiver))
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, err
		}
	} else {
		return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, errors.NewMsg("Invalid Transaction type")
	}
	return state, nil
}

func ExecuteTransactions(txs []Transaction, state mpt.MerklePatriciaTrie) (mpt.MerklePatriciaTrie, error) {
	for _, tx := range txs {
		var err error
		state, err = ExecuteTransaction(tx, state)
		if err != nil {
			return mpt.MerklePatriciaTrie{Root: common.EmptyHash}, err
		}
	}
	return state, nil
}

func CreateCoinbaseTransaction(miner common.Address) Transaction {
	tx := Transaction{
		TxType:          2,
		SenderPublicKey: common.EmptyPublicKey,
		SenderSignature: common.EmptySignature,
		Receiver:        miner,
		Value:           COINBASE_VALUE,
		Fee:             0,
		Nonce:           0,
	}
	return tx
}
