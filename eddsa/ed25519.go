package eddsa

import (
	"crypto/rand"

	"golang.org/x/crypto/ed25519"
)

type curveED25519 struct {
	EdDSA
}

// NewCurveED25519 creates a new Edwards-Curve Digital Signature Algorithm (EdDSA) instance
func NewCurveED25519() EdDSA {
	return &curveED25519{}
}

// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func (e *curveED25519) GenerateKeyPair() ([]byte, []byte, error) {
	return ed25519.GenerateKey(rand.Reader)
}

// Sign signs the message with privateKey and returns a signature. It will
// panic if len(privateKey) is not PrivateKeySize = 64.
func (e *curveED25519) Sign(privateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}
