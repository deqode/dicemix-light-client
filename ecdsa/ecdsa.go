package ecdsa

// ECDSA - The main interface P256 curve.
type ECDSA interface {
	GenerateKeyPair() ([]byte, []byte, error)
	Sign([]byte, []byte) []byte
}
