package rng

// RNG - The main interface chacha20 DiceMixRng.
type RNG interface {
	GetFieldElement(dicemix DiceMixRng) uint64
	GetBytes(dicemix DiceMixRng, bytes uint8) []byte
}
