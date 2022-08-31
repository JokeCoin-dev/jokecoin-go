package common

import (
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

func ToHash(b []byte) Hash {
	return *(*Hash)(b)
}

func PublicKeyToAddress(pk ed25519.PublicKey) Address {
	return sha3.Sum256(pk[:])
}
