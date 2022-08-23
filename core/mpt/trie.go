package mpt

import (
	"jokecoin-go/core/common"
	"jokecoin-go/core/database"
	"jokecoin-go/core/errors"
	"jokecoin-go/core/utils"
)

type MerklePatriciaTrie struct {
	Root common.Hash
}

func NewTrie() MerklePatriciaTrie {
	return MerklePatriciaTrie{
		Root: common.EmptyHash,
	}
}

// Get returns the value for the given key.
// Return (nil,nil) if the key does not exist.
func (t MerklePatriciaTrie) Get(key []byte) ([]byte, error) {
	db := database.GetDB()
	node := t.Root
	nibbles := FromBytes(key)
	for {
		if node == common.EmptyHash {
			return nil, nil
		}
		b, err := db.Get(node[:])
		if err != nil {
			return nil, errors.Wrap(err, errors.DatabaseError)
		}
		var inode Node
		err = utils.Deserialize(b, &inode)
		if err != nil {
			return nil, errors.Wrap(err, errors.DataConsistencyError)
		}
		if leaf, ok := inode.(LeafNode); ok {
			matched := PrefixMatchedLen(nibbles, leaf.Path)
			if matched != len(leaf.Path) || matched != len(nibbles) {
				return nil, nil
			}
			return leaf.Value, nil
		} else if branch, ok := inode.(BranchNode); ok {
			if len(nibbles) == 0 {
				return nil, errors.New(errors.DataConsistencyError)
			}
			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			node = branch.Branches[b]
		} else if ext, ok := inode.(ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, nibbles)
			if matched < len(ext.Path) {
				return nil, errors.New(errors.DataConsistencyError)
			}
			nibbles = nibbles[matched:]
			node = ext.Next
		}
	}
}

func (t MerklePatriciaTrie) Put(key []byte, value []byte) (MerklePatriciaTrie, error) {
	nibble := FromBytes(key)
	n, err := put(t.Root, nibble, value)
	if err != nil {
		return MerklePatriciaTrie{}, err
	}
	return MerklePatriciaTrie{Root: n}, nil
}

func put(node common.Hash, key []Nibble, value []byte) (common.Hash, error) {
	nibbles := key
	db := database.GetDB()
	if node == common.EmptyHash {
		leaf := LeafNode{Path: nibbles, Value: value}
		return leaf.WriteDB()
	} else {
		b, err := db.Get(node[:])
		if err != nil {
			return common.EmptyHash, errors.Wrap(err, errors.DatabaseError)
		}
		var inode Node
		err = utils.Deserialize(b, &inode)
		if err != nil {
			return common.EmptyHash, errors.Wrap(err, errors.DataConsistencyError)
		}
		if leaf, ok := inode.(LeafNode); ok {
			if len(leaf.Path) != len(nibbles) {
				return common.EmptyHash, errors.New(errors.DataConsistencyError)
			}
			matched := PrefixMatchedLen(leaf.Path, nibbles)
			if matched == len(nibbles) {
				leaf := LeafNode{Path: nibbles, Value: value}
				return leaf.WriteDB()
			} else {
				branch := NewBranch()
				leaf1 := LeafNode{leaf.Path[matched+1:], leaf.Value}
				hash1, err := leaf1.WriteDB()
				if err != nil {
					return common.EmptyHash, err
				}
				branch.Branches[leaf.Path[matched]] = hash1
				leaf2 := LeafNode{nibbles[matched+1:], value}
				hash2, err := leaf2.WriteDB()
				if err != nil {
					return common.EmptyHash, err
				}
				branch.Branches[nibbles[matched]] = hash2
				hashBranch, err := branch.WriteDB()
				if err != nil {
					return common.EmptyHash, err
				}
				if matched > 0 {
					ext := ExtensionNode{leaf.Path[:matched], hashBranch}
					return ext.WriteDB()
				} else {
					return hashBranch, nil
				}
			}
		} else if branch, ok := inode.(BranchNode); ok {
			b, remaining := nibbles[0], nibbles[1:]
			newNode, err := put(branch.Branches[b], remaining, value)
			if err != nil {
				return common.EmptyHash, err
			}
			branch.Branches[b] = newNode
			return branch.WriteDB()
		} else if ext, ok := inode.(ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, nibbles)
			if len(ext.Path) >= len(nibbles) {
				return common.EmptyHash, errors.New(errors.DataConsistencyError)
			}
			if matched < len(ext.Path) {
				extNibbles, branchNibble, extRemainingnibbles := ext.Path[:matched], ext.Path[matched], ext.Path[matched+1:]
				nodeBranchNibble, nodeLeafNibbles := nibbles[matched], nibbles[matched+1:]
				branch := NewBranch()
				if len(extRemainingnibbles) == 0 {
					branch.Branches[branchNibble] = ext.Next
				} else {
					newExt := ExtensionNode{extRemainingnibbles, ext.Next}
					hash, err := newExt.WriteDB()
					if err != nil {
						return common.EmptyHash, err
					}
					branch.Branches[branchNibble] = hash
				}
				leaf := LeafNode{nodeLeafNibbles, value}
				enodeLeaf, err := leaf.WriteDB()
				if err != nil {
					return common.EmptyHash, err
				}
				branch.Branches[nodeBranchNibble] = enodeLeaf
				hashBranch, err := branch.WriteDB()
				if err != nil {
					return common.EmptyHash, err
				}
				if len(extNibbles) == 0 {
					return hashBranch, nil
				} else {
					newExt := ExtensionNode{extNibbles, hashBranch}
					return newExt.WriteDB()
				}
			} else {
				newNext, err := put(ext.Next, nibbles[matched:], value)
				if err != nil {
					return common.EmptyHash, err
				}
				ext.Next = newNext
				return ext.WriteDB()
			}
		} else {
			return common.EmptyHash, errors.New(errors.DataConsistencyError)
		}
	}
}
