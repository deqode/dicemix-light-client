package ecdh

import (
	"crypto"
)

// ECDH - The main interface ECDH.
type ECDH interface {
	GenerateKeyPair() (crypto.PrivateKey, crypto.PublicKey, error)
	Marshal(p crypto.PublicKey) []byte
	Unmarshal(data []byte) (crypto.PublicKey, bool)
	GenerateSharedSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error)
}
