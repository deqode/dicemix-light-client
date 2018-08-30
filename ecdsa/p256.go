package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
)

type curveP256 struct {
	ECDSA
}

// NewCurveECDSA creates a new Elliptic Curve Digital Signature Algorithm  instance
func NewCurveECDSA() ECDSA {
	return &curveP256{}
}

// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func (e *curveP256) GenerateKeyPair() ([]byte, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return nil, nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)

	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return nil, nil, err
	}

	return publicKeyBytes, privateKeyBytes, nil
}

// Sign signs the message with privateKey and returns a signature. It will
// return nil if error occurs
func (e *curveP256) Sign(privateKeyBytes, message []byte) []byte {
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return nil
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message)
	if err != nil {
		return nil
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return signature
}
