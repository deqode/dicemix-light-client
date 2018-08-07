package nike

import (
	"fmt"
	"log"
	"sync"

	"../ecdh"
	"../rng"
	"../utils"
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
func (n *nike) GenerateKeys(state *utils.State) {
	n.Lock()
	defer n.Unlock()
	// generate random key pair
	ecdh := ecdh.NewCurve25519ECDH()
	var err error
	(*state).Kesk, (*state).Kepk, err = ecdh.GenerateKeyPair()

	if err != nil {
		log.Fatalf("Error: generating NIKE key pair %v", err)
	}

	fmt.Printf("MY KESK - %v\n\n", (*state).Kesk)
	fmt.Printf("MY KEPK - %v\n\n", (*state).Kepk)
}

// DeriveSharedKeys - derives shared keys for all peers
// generates RNG based on shared key using ChaCha20
func (n *nike) DeriveSharedKeys(state *utils.State) {
	n.Lock()
	defer n.Unlock()

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

	fmt.Printf("\n--------------------------------------\n")

	for _, peer := range state.Peers {
		fmt.Printf("ID - %v\n", peer.ID)
		fmt.Printf("PK - %v\n", peer.PubKey)
		fmt.Printf("Shared Key - %v\n\n", peer.SharedKey)
	}

	fmt.Printf("\n--------------------------------------\n")
}
