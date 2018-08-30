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
	// generates random privateKey
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// check for error
	if err != nil {
		return nil, nil, err
	}

	// obtain publicKey from privateKey and Marshal into bytes
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	// Marshal privateKey into bytes
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)

	// check for error
	if err != nil {
		return nil, nil, err
	}

	return publicKeyBytes, privateKeyBytes, nil
}

// Sign signs the message with privateKey and returns a signature. It will
// return nil if error occurs
func (e *curveP256) Sign(privateKeyBytes, message []byte) []byte {
	// get privateKey object from privateKeyBytes
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)

	// check for error
	if err != nil {
		return nil
	}

	// sign message with privateKey
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message)

	if err != nil {
		return nil
	}

	// obtain signature from r, s of type *big.Int
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return signature
}
