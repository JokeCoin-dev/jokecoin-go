package common

import (
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/utils"
)

type AccountData struct {
	Balance     int64
	Nonce       int64
	CodeHash    Hash
	StorageHash Hash
}

func (a AccountData) GetHash() Hash {
	return sha3.Sum256(utils.MustSerialize(a))
}
