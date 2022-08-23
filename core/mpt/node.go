package mpt

import (
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/common"
	"jokecoin-go/core/database"
	"jokecoin-go/core/errors"
	"jokecoin-go/core/utils"
)

type NodeType byte

const (
	TypeLeaf      = 1
	TypeBranch    = 2
	TypeExtension = 3
)

const BranchSize = 16

type Node struct {
	Type     NodeType
	Path     []Nibble      // Leaf & Extension
	Value    []byte        //Leaf
	Branches []common.Hash //Branch
	Next     common.Hash   //Extension
}

/*func (n EncodedNode) GetHash() common.Hash {
	b, err := utils.Serialize(n)
	utils.PanicIfErr(err)
	return sha3.Sum256(b)
}

func (n EncodedNode) WriteDB() error {
	db := database.GetDB()
	hash := n.GetHash()
	err := db.Put(hash[:], utils.MustSerialize(n), nil)
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	} else {
		return nil
	}
}

func CreateENode(child interface{}) (common.Hash, error) {
	var enode EncodedNode
	if n, ok := child.(LeafNode); ok {
		enode = EncodedNode{Type: TypeLeaf, Data: utils.MustSerialize(n)}
	} else if n, ok := child.(BranchNode); ok {
		enode = EncodedNode{Type: TypeBranch, Data: utils.MustSerialize(n)}
	} else if n, ok := child.(ExtensionNode); ok {
		enode = EncodedNode{Type: TypeExtension, Data: utils.MustSerialize(n)}
	} else {
		log.Panic("invalid node type")
	}
	err := enode.WriteDB()
	if err != nil {
		return common.EmptyHash, errors.Wrap(err, errors.DatabaseError)
	}
	return enode.GetHash(), nil
}*/

func (n Node) GetHash() common.Hash {
	b, err := utils.Serialize(n)
	utils.PanicIfErr(err)
	return sha3.Sum256(b)
}

func (n Node) WriteDB() (common.Hash, error) {
	db := database.GetDB()
	hash := n.GetHash()
	err := db.Put(hash[:], utils.MustSerialize(n))
	if err != nil {
		return common.EmptyHash, errors.Wrap(err, errors.DatabaseError)
	} else {
		return hash, nil
	}
}
func NewBranch() *Node {
	b := &Node{Type: TypeBranch}
	b.Branches = make([]common.Hash, BranchSize)
	for i := 0; i < BranchSize; i++ {
		b.Branches[i] = common.EmptyHash
	}
	return b
}

func NewLeaf(path []Nibble, value []byte) *Node {
	n := &Node{Type: TypeLeaf, Path: path, Value: value}
	return n
}

func NewExtension(path []Nibble, next common.Hash) *Node {
	n := &Node{Type: TypeExtension, Path: path, Next: next}
	return n
}
