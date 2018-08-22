package eddsa

// EdDSA - The main interface ed25519.
type EdDSA interface {
	GenerateKeyPair() ([]byte, []byte, error)
	Sign([]byte, []byte) []byte
}
