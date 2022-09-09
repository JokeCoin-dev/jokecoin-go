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
var EmptyPublicKey ed25519.PublicKey
var EmptySignature Signature = [SignatureLen]byte{}

type Hash [HashLen]byte
type Address [AddressLen]byte
type Signature [SignatureLen]byte

func init() {
	EmptyHash = sha3.Sum256([]byte(""))
	EmptyPublicKey = make([]byte, PublicKeyLen)
}
