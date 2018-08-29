package utils

import (
	"crypto"
	"dicemix_client/rng"
	"math/rand"
	"time"

	base58 "github.com/jbenet/go-base58"
)

const (
	// MaxAllowedMessages - Basic sanity check to avoid weird inputs
	// if messages are more than MaxAllowedMessages then close connection
	MaxAllowedMessages = 1000

	// ResponseWait - Time to wait for response from server
	// if server does'nt response within ResponseWait seconds then close connection
	ResponseWait = 30
)

// Peers - Stores all Peers Info
type Peers struct {
	ID             int32
	PubKey         []byte
	NumMsgs        uint32
	SharedKey      []byte
	Dicemix        rng.DiceMixRng
	DCSimpleVector [][]byte
	Ok             bool
	Confirmation   []byte
}

// Session stores information of current Session
type session struct {
	Ltsk      []byte
	Ltpk      []byte
	SessionID uint64
	MyID      int32
	Kesk      crypto.PrivateKey
	NextKesk  crypto.PrivateKey
	Kepk      crypto.PublicKey
	NextKepk  crypto.PublicKey
}

// State - stores state info for current run
// TODO: remove ltsk and ltpk from state and
// store them in more persistent storage
type State struct {
	Session        session
	Peers          []Peers
	AllMsgHashes   []uint64
	MyDC           []uint64
	MyOk           bool
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
