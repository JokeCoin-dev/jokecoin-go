package mpt

import (
	"encoding/gob"
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

type LeafNode struct {
	Path  []Nibble
	Value []byte
}

type BranchNode struct {
	Branches [BranchSize]common.Hash
}

type ExtensionNode struct {
	Path []Nibble
	Next common.Hash
}

type Node interface {
	GetHash() common.Hash
	WriteDB() (common.Hash, error)
}

func init() {
	gob.Register(LeafNode{})
	gob.Register(BranchNode{})
	gob.Register(ExtensionNode{})
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

func (n LeafNode) GetHash() common.Hash {
	b, err := utils.Serialize(n)
	utils.PanicIfErr(err)
	return sha3.Sum256(b)
}

func (n LeafNode) WriteDB() (common.Hash, error) {
	db := database.GetDB()
	hash := n.GetHash()
	t := Node(n)
	err := db.Put(hash[:], utils.MustSerialize(&t))
	if err != nil {
		return common.EmptyHash, errors.Wrap(err, errors.DatabaseError)
	} else {
		return hash, nil
	}
}
func NewBranch() BranchNode {
	b := BranchNode{}
	for i := 0; i < BranchSize; i++ {
		b.Branches[i] = common.EmptyHash
	}
	return b
}

func (n BranchNode) GetHash() common.Hash {
	b, err := utils.Serialize(n)
	utils.PanicIfErr(err)
	return sha3.Sum256(b)
}

func (n BranchNode) WriteDB() (common.Hash, error) {
	db := database.GetDB()
	hash := n.GetHash()
	t := Node(n)
	err := db.Put(hash[:], utils.MustSerialize(&t))
	if err != nil {
		return common.EmptyHash, errors.Wrap(err, errors.DatabaseError)
	} else {
		return hash, nil
	}
}

func (n ExtensionNode) GetHash() common.Hash {
	b, err := utils.Serialize(n)
	utils.PanicIfErr(err)
	return sha3.Sum256(b)
}

func (n ExtensionNode) WriteDB() (common.Hash, error) {
	db := database.GetDB()
	hash := n.GetHash()
	t := Node(n)
	err := db.Put(hash[:], utils.MustSerialize(&t))
	if err != nil {
		return common.EmptyHash, errors.Wrap(err, errors.DatabaseError)
	} else {
		return hash, nil
	}
}
