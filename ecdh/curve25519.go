package ecdh

import (
	"crypto"
	"crypto/rand"
	"sync"

	ecdh "github.com/wsddn/go-ecdh"
)

type curve25519ECDH struct {
	ECDH
	sync.Mutex
}

// NewCurve25519ECDH creates a new ECDH instance that uses djb's curve25519
// elliptical curve.
func NewCurve25519ECDH() ECDH {
	return &curve25519ECDH{}
}

// GenerateKeyPair creates new PrivateKey and PublicKey that uses djb's curve25519
// elliptical curve.
func (e *curve25519ECDH) GenerateKeyPair() (crypto.PrivateKey, crypto.PublicKey, error) {
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.GenerateKey(rand.Reader)
}

// Marshal converts crypto.PublicKey into byte[]
func (e *curve25519ECDH) Marshal(p crypto.PublicKey) []byte {
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.Marshal(p)
}

// Unmarshal converts byte[] to crypto.PublicKey
func (e *curve25519ECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.Unmarshal(data)
}

// Marshal converts crypto.PrivateKey into byte[]
// used in KESK stage
func (e *curve25519ECDH) MarshalSK(p crypto.PrivateKey) []byte {
	pri := p.(*[32]byte)
	return pri[:]
}

// GenerateSharedSecret creates shared key using our private key and others public key
func (e *curve25519ECDH) GenerateSharedSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.GenerateSharedSecret(privKey, pubKey)
}

// GenerateConfirmation creates confirmation on messages
// NOTE: this is dummy implementation only for testing purpose
func (e *curve25519ECDH) GenerateConfirmation(messages [][]byte, privKey crypto.PrivateKey) ([]byte, error) {
	// TODO: create proper function to generate confirmations on messages
	return messages[0], nil
}
