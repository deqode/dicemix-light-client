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
	e.Lock()
	defer e.Unlock()
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.GenerateKey(rand.Reader)
}

// Marshal converts crypto.PublicKey into byte[]
func (e *curve25519ECDH) Marshal(p crypto.PublicKey) []byte {
	e.Lock()
	defer e.Unlock()
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.Marshal(p)
}

// Unmarshal converts byte[] to crypto.PublicKey
func (e *curve25519ECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	e.Lock()
	defer e.Unlock()
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.Unmarshal(data)
}

// GenerateSharedSecret creates shared key using our private key and others public key
func (e *curve25519ECDH) GenerateSharedSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	e.Lock()
	defer e.Unlock()
	var ecdhCurve = ecdh.NewCurve25519ECDH()
	return ecdhCurve.GenerateSharedSecret(privKey, pubKey)
}
