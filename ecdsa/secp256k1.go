package ecdsa

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type curveS256 struct {
	ECDSA
}

// NewCurveECDSA creates a new Elliptic Curve Digital Signature Algorithm  instance
func NewCurveECDSA() ECDSA {
	return &curveS256{}
}

// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func (e *curveS256) GenerateKeyPair() ([]byte, []byte, error) {
	// generates random privateKey
	privateKey, err := btcec.NewPrivateKey(btcec.S256())

	// check for error
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes := privateKey.PubKey().SerializeCompressed()
	privateKeyBytes := privateKey.Serialize()

	return publicKeyBytes, privateKeyBytes, nil
}

// Sign signs the message with privateKey and returns a signature. It will
// return nil if error occurs
func (e *curveS256) Sign(privateKeyBytes, message []byte) []byte {
	// obtain private key object from bytes
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)

	// hash message that is to be sent to server
	messageHash := chainhash.DoubleHashB(message)

	// sign message hash using private key
	signature, err := privateKey.Sign(messageHash)

	// if error occured return nil
	if err != nil {
		return nil
	}

	// serialize and return the signature.
	return signature.Serialize()
}
