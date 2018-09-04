package dc

import (
	"dicemix-client/field"
	"dicemix-client/utils"

	log "github.com/sirupsen/logrus"
)

// to expose DC-NET methods
type dcNet struct {
	DC
}

// NewDCNetwork creates a new DC instance
func NewDCNetwork() DC {
	return &dcNet{}
}

// RunDCSimple - Runs DC-Simple with slot reservation
func (d *dcNet) RunDCSimple(state *utils.State) {
	// initaializing variables
	slots := make([]int, state.MyMsgCount)
	peersCount := uint32(len(state.Peers))
	var i, j uint32
	totalMsgsCount := messageCount(state.MyMsgCount, state.Peers)

	// insanity check
	if totalMsgsCount > utils.MaxAllowedMessages {
		log.Fatalf("Limit Exceeded: More than %d messages in tx", utils.MaxAllowedMessages)
	}

	slots, state.MyOk = obtainSlots(state, totalMsgsCount)

	// if one of my message hashes (root) is missing
	if !state.MyOk {
		// Even though the run will be aborted (because we send my_ok = false), transmit the
		// message in a deterministic slot. This enables the peers to recompute our commitment.
		for i = 0; i < state.MyMsgCount; i++ {
			slots[i] = int(i)
		}
	}

	// array of |totalMsgsCount| arrays of slot_size bytes, all initalized with 0
	state.DCSimpleVector = make([][]byte, totalMsgsCount)

	// reserve 20 bytes (160 bits) for each slot
	// to store messages of ours and peers
	for j = 0; j < totalMsgsCount; j++ {
		state.DCSimpleVector[j] = make([]byte, 20)
	}

	// store our all messages (byte encoded) in slot reserved
	for j = 0; j < state.MyMsgCount; j++ {
		state.DCSimpleVector[slots[j]] = utils.Base58StringToBytes(state.MyMessages[j])
	}

	log.Info("Slot's = ", state.DCSimpleVector)

	// encode messages in slots
	for i = 0; i < peersCount; i++ {
		for j = 0; j < totalMsgsCount; j++ {
			// xor operation - dc_simple_vector[j] = dc_simple_vector[j] + <randomness for chacha20>
			xorBytes(state.DCSimpleVector[j], state.DCSimpleVector[j], state.Peers[i].Dicemix.GetBytes(20))
		}
	}

	log.Info("My DC-SIMPLE vector = ", state.DCSimpleVector)
}

// Run a DC-net with exponential encoding
// generates my_dc[]
func (d *dcNet) DeriveMyDCVector(state *utils.State) {
	// initialize variables
	var i, j uint32
	peersCount := uint32(len(state.Peers))
	totalMsgsCount := messageCount(state.MyMsgCount, state.Peers)
	state.MyDC = make([]uint64, totalMsgsCount)

	// generates power sums of message_hashes
	// my_dc[i] := my_dc[i] (+) (my_msg_hashes[j] ** (i + 1))
	for j = 0; j < state.MyMsgCount; j++ {
		// generates 64 bit hash of my_message[j]
		state.MyMessagesHash[j] = shortHash(state.MyMessages[j])
		var pow uint64 = 1
		for i = 0; i < totalMsgsCount; i++ {
			pow = power(state.MyMessagesHash[j], pow)
			state.MyDC[i] = field.NewField(state.MyDC[i]).Add(field.NewField(pow)).Value()
		}
	}

	// encode power sums
	// my_dc[i] := my_dc[i] (+) (sgn(my_id - p.id) (*) p.dicemix.get_field_element())
	for j = 0; j < peersCount; j++ {
		for i = 0; i < totalMsgsCount; i++ {
			var op2 = field.NewField(state.Peers[j].Dicemix.GetFieldElement())
			if state.Session.MyID < state.Peers[j].ID {
				op2 = op2.Neg()
			}
			state.MyDC[i] = field.NewField(state.MyDC[i]).Add(op2).Value()
		}
	}

	log.Info("My Msg Hashes = ", state.MyMessagesHash)
	log.Info("My DC-EXP vector = ", state.MyDC)
}

// Verify that every peer agrees to proceed
func (d *dcNet) VerifyProceed(state *utils.State) bool {
	totalMsgsCount := messageCount(state.MyMsgCount, state.Peers)
	slots, ok := obtainSlots(state, totalMsgsCount)

	if !ok {
		return false
	}

	for _, index := range slots {
		s := shortHash(utils.BytesToBase58String(state.AllMessages[index]))
		if state.AllMsgHashes[index] != reduce(s) {
			return false
		}
	}

	return true
}
