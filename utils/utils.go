package utils

import (
	"crypto"
	"math/rand"
	"time"

	"../rng"
	base58 "github.com/jbenet/go-base58"
)

const (
	// MaxAllowedMessages - Basic sanity check to avoid weird inputs
	MaxAllowedMessages = 1000

	// ResponseWait - Time to wait for response from server.
	ResponseWait = 30
)

// Peers - Stores all Peers Info
type Peers struct {
	ID             int32
	PubKey         []byte
	NextPubKey     []byte
	NumMsgs        uint32
	SharedKey      []byte
	Dicemix        rng.DiceMixRng
	DCSimpleVector [][]byte
	Ok             bool
	Confirmation   []byte
}

// State - stores state info for current run
type State struct {
	Ltsk           []byte
	Ltpk           []byte
	Run            int
	SessionID      uint64
	Peers          []Peers
	TotalMsgsCount uint32
	AllMsgHashes   []uint64
	MyID           int32
	MyDC           []uint64
	MyOk           bool
	Kesk           crypto.PrivateKey
	NextKesk       crypto.PrivateKey
	Kepk           crypto.PublicKey
	NextKepk       crypto.PublicKey
	MyMessages     []string
	MyMessagesHash []uint64
	MyMsgCount     uint32
	DCSimpleVector [][]byte
	AllMessages    [][]byte
}

// GenerateMessage - generates a random 20 byte string (160 bits)
// returns string encoded with Base58 format
func GenerateMessage() string {
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, 20)
	rand.Read(token)
	return BytesToBase58String(token)
}

// BytesToBase58String - converts []byte to Base58 Encoded string
func BytesToBase58String(bytes []byte) string {
	return base58.Encode(bytes)
}

// Base58StringToBytes - converts Base58 Encoded string to []byte
func Base58StringToBytes(str string) []byte {
	return base58.Decode(str)
}
