package rng

import (
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"

	"github.com/codahale/chacha20"
	log "github.com/sirupsen/logrus"
)

// DiceMixRng -- data structure to hold stream
// and random numbers generated through chacha20
type DiceMixRng struct {
	chachaStream    cipher.Stream
	chachaStreamErr error
	chachaExpRng    []byte
}

// NewRng -- creates DiceMixRng object using seed provided
func NewRng(seed []byte) DiceMixRng {
	// default -- using nonce value as 0
	nonceHex := "0000000000000000"
	nonce, err := hex.DecodeString(nonceHex)

	if err != nil {
		log.Fatal("Error Occured: ", err)
	}

	var dicemix = DiceMixRng{}
	dicemix.chachaStream, dicemix.chachaStreamErr = chacha20.New(seed, nonce)

	if dicemix.chachaStreamErr != nil {
		log.Fatal("Error Occured: ", dicemix.chachaStreamErr)
	}

	// generate random number for DC Exponential by extracting first 8 bytes
	dicemix.chachaExpRng = getPRG(dicemix.chachaStream, 8)

	return dicemix
}

// GetFieldElement - converts []byte of len 8 to uint64
func (d *DiceMixRng) GetFieldElement() uint64 {
	return uint64(binary.LittleEndian.Uint64(d.chachaExpRng))
}

// GetBytes - returns 20 byte[]
func (d *DiceMixRng) GetBytes(bytes uint8) []byte {
	return getPRG(d.chachaStream, bytes)
}

// generates Rng in form of bytes[] and string from provided stream
func getPRG(stream cipher.Stream, pos uint8) []byte {
	src := make([]byte, 64)
	dst := make([]byte, 64)

	// stores stream bytes into dst[]
	stream.XORKeyStream(dst, src)

	// extracts first |pos| bytes from dst[]
	return dst[:pos]
}
