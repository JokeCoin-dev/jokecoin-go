package block

import (
	"crypto/ed25519"
	"jokecoin-go/core/common"
	"jokecoin-go/core/database"
	"jokecoin-go/core/mpt"
	"jokecoin-go/core/utils"
	"testing"
)

func TestExecuteTransaction(t *testing.T) {
	utils.PanicIfErr(database.InitLevelDB("test_leveldb"))
	pubkey1, privkey1, err := ed25519.GenerateKey(nil)
	utils.PanicIfErr(err)
	addr1 := common.PublicKeyToAddress(pubkey1)
	pubkey2, _, err := ed25519.GenerateKey(nil)
	utils.PanicIfErr(err)
	addr2 := common.PublicKeyToAddress(pubkey2)
	TX1 := Transaction{
		TxType:   2,
		Receiver: addr1,
		Value:    COINBASE_VALUE,
	}
	TX2 := Transaction{
		TxType:          1,
		SenderPublicKey: pubkey1,
		Receiver:        addr2,
		Value:           100,
		Fee:             0,
		Nonce:           1,
	}
	TX2.Sign(privkey1)
	root := mpt.MerklePatriciaTrie{Root: common.EmptyHash}
	TXs := []Transaction{TX1, TX2}
	root, err = ExecuteTransactions(TXs, root)
	utils.PanicIfErr(err)
	b, err := root.Get(addr1[:])
	utils.PanicIfErr(err)
	account1 := common.AccountData{}
	utils.MustDeserialize(b, &account1)
	b, err = root.Get(addr2[:])
	utils.PanicIfErr(err)
	account2 := common.AccountData{}
	utils.MustDeserialize(b, &account2)
	utils.Assert(account1.Balance == COINBASE_VALUE-100)
	utils.Assert(account1.Nonce == 1)
	utils.Assert(account2.Balance == 100)
	utils.PanicIfErr(database.GetDB().Close())
}
