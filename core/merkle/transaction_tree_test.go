package merkle

import (
	"crypto/ed25519"
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/block"
	"jokecoin-go/core/common"
	"jokecoin-go/core/database"
	"jokecoin-go/core/utils"
	"testing"
)

func TestBuildTransactionTree(t *testing.T) {
	utils.PanicIfErr(database.InitLevelDB("test_leveldb"))
	TX := block.Transaction{
		TxType:          0,
		SenderPublicKey: ed25519.PublicKey{},
		SenderSignature: common.Signature{},
		Receiver:        common.Address{},
		Value:           0,
		GasLimit:        0,
		Fee:             0,
		Nonce:           0,
		Data:            []byte("A Transaction"),
	}
	var TXs []block.Transaction
	for i := 0; i < 10; i++ {
		TX.Nonce = int64(i)
		TXs = append(TXs, TX)
	}
	root := BuildTransactionTree(TXs)
	t.Log(root)
	db := database.GetDB()
	x := &Node{}
	h := TXs[0].ComputeHash()
	h2 := sha3.Sum256(h[:])
	for i := 1; i <= 5; i++ {
		b, err := db.Get(h2[:])
		if err != nil {
			t.Fatal(err)
		}
		err = utils.Deserialize(b, x)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(x)
		h2 = x.Parent
	}
	utils.Assert(h2 == EmptyTree)
	utils.PanicIfErr(db.Close())
}
