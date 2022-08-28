package common

import (
	"crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

const HashLen = 32
const AddressLen = HashLen
const PublicKeyLen = ed25519.PublicKeySize
const PrivateKeyLen = ed25519.PrivateKeySize
const SignatureLen = ed25519.SignatureSize

var EmptyHash Hash

type Hash [HashLen]byte
type Address [AddressLen]byte
type PublicKey [PublicKeyLen]byte
type PrivateKey [PrivateKeyLen]byte
type Signature [SignatureLen]byte

func init() {
	EmptyHash = sha3.Sum256([]byte(""))
}

func ToHash(b []byte) Hash {
	return *(*Hash)(b)
}
