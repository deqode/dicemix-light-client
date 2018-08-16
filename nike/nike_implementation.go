package nike

import (
	"sync"

	"../ecdh"
	"../rng"
	"../utils"
	log "github.com/sirupsen/logrus"
)

type nike struct {
	NIKE
	sync.Mutex
}

// NewNike creates a new NIKE instance
func NewNike() NIKE {
	return &nike{}
}

// KeyExchange -- generates random NIKE keypair, message
// broadcasts self public-key
// receive other peers public-keys[]
// mode = 0 to generate (my_kesk, my_kepk)
// mode = 1 to generate (my_next_kesk, my_next_kepk)
func (n *nike) GenerateKeys(state *utils.State, mode int) {
	// generate random key pair
	ecdh := ecdh.NewCurve25519ECDH()
	var err error
	if mode == 0 {
		(*state).Kesk, (*state).Kepk, err = ecdh.GenerateKeyPair()
	} else if mode == 1 {
		(*state).NextKesk, (*state).NextKepk, err = ecdh.GenerateKeyPair()
	}

	if err != nil {
		log.Fatalf("Error: generating NIKE key pair %v", err)
	}
}

// DeriveSharedKeys - derives shared keys for all peers
// generates RNG based on shared key using ChaCha20
func (n *nike) DeriveSharedKeys(state *utils.State) {
	ecdh := ecdh.NewCurve25519ECDH()
	peersCount := len((*state).Peers)
	for i := 0; i < peersCount; i++ {
		var pubkey, res = ecdh.Unmarshal((*state).Peers[i].PubKey)
		if !res {
			log.Fatalf("Error: generating NIKE Shared Keys %v", res)
		}
		var err error
		(*state).Peers[i].SharedKey, err = ecdh.GenerateSharedSecret((*state).Kesk, pubkey)

		if err != nil {
			log.Fatalf("Error: generating NIKE Shared Keys %v", err)
		}

		(*state).Peers[i].Dicemix = rng.NewRng((*state).Peers[i].SharedKey)
	}
}
