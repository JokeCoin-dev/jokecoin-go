package merkle

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/block"
	"jokecoin-go/core/common"
	"jokecoin-go/core/database"
	"jokecoin-go/core/utils"
)

type Node struct {
	Parent     common.Hash
	LeftChild  common.Hash
	RightChild common.Hash
}

var EmptyTree common.Hash

func (n *Node) Encode() []byte {
	b, err := utils.Serialize(n)
	utils.PanicIfErr(err)
	return b
}

func DecodeNode(b []byte) Node {
	n := &Node{}
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(n)
	utils.PanicIfErr(err)
	return *n
}

func init() {
	EmptyTree = common.EmptyHash
}

func BuildTransactionTree(TXs []block.Transaction) common.Hash {
	root := buildTransactionTree1(TXs)
	buildTransactionTree2(root, EmptyTree)
	return root
}

func buildTransactionTree1(TXs []block.Transaction) common.Hash {
	db := database.GetDB()
	if len(TXs) == 0 {
		return EmptyTree
	}
	if len(TXs) == 1 {
		a := TXs[0].ComputeHash()
		hash := sha3.Sum256(a[:])
		now := Node{LeftChild: EmptyTree, RightChild: EmptyTree, Parent: EmptyTree}
		db.Put(hash[:], now.Encode())
		return hash
	}
	mid := (len(TXs) + 1) / 2
	leftChild := buildTransactionTree1(TXs[:mid])
	rightChild := buildTransactionTree1(TXs[mid:])
	hash := sha3.Sum256(append(leftChild[:], rightChild[:]...))
	now := Node{LeftChild: leftChild, RightChild: rightChild, Parent: EmptyTree}
	db.Put(hash[:], now.Encode())
	return hash
}
func buildTransactionTree2(now common.Hash, fa common.Hash) {
	db := database.GetDB()
	b, err := db.Get(now[:])
	utils.PanicIfErr(err)
	n := DecodeNode(b)
	n.Parent = fa
	db.Put(now[:], n.Encode())
	if n.LeftChild != EmptyTree {
		buildTransactionTree2(n.LeftChild, now)
		buildTransactionTree2(n.RightChild, now)
	}
}
