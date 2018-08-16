package dc

import (
	"os"

	"../field"
	"../utils"
	"github.com/shomali11/util/xhashes"
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
	state.MyOk = true
	var i, j uint32
	totalMsgsCount := state.MyMsgCount

	for _, peer := range state.Peers {
		totalMsgsCount += peer.NumMsgs
	}

	if totalMsgsCount > 1000 {
		log.Fatal("Limit Exceeded: More than 1000 messages in tx")
		os.Exit(1)
	}

	for i := range slots {
		slots[i] = -1
	}

	// Run an ordinary DC-net with slot reservations
	for j = 0; j < state.MyMsgCount; j++ {
		index, count := -1, 0
		for i = 0; i < totalMsgsCount; i++ {
			if state.AllMsgHashes[i] == reduce(state.MyMessagesHash[j]) {
				index, count = int(i), int(count+1)
			}
		}

		// if there is exactly one i
		// with all_msg_hashes[i] = my_msg_hashes[j] then
		if count == 1 {
			slots[j] = index
		} else {
			state.MyOk = false
		}
	}

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

	for i = 0; i < peersCount; i++ {
		for j = 0; j < totalMsgsCount; j++ {
			// encode messages in slots
			// xor operation - dc_simple_vector[j] = dc_simple_vector[j] + <randomness for chacha20>
			xorBytes(state.DCSimpleVector[j], state.DCSimpleVector[j], state.Peers[i].Dicemix.GetBytes(20))
		}
	}

	log.Info("My DC-SIMPLE vector = ", state.DCSimpleVector)
}

// Resolve the DC-net
func (d *dcNet) ResolveDCNet(state *utils.State) {
	var i, j uint32
	peersCount := uint32(len(state.Peers))
	totalMsgsCount := state.MyMsgCount

	for _, peer := range state.Peers {
		totalMsgsCount += peer.NumMsgs
	}
	state.AllMessages = state.DCSimpleVector

	for i = 0; i < peersCount; i++ {
		for j = 0; j < totalMsgsCount; j++ {
			// decodes messages from slots by cancelling out randomness introduced in DC-Simple
			// xor operation - all_messages[j] = dc_simple_vector[j] + <randomness for chacha20>
			xorBytes(state.AllMessages[j], state.AllMessages[j], state.Peers[i].DCSimpleVector[j])
		}
	}

	log.Info("Resolved DC-NET vector = ", state.AllMessages)
}

// Run a DC-net with exponential encoding
// generates my_dc[]
func (d *dcNet) DeriveMyDCVector(state *utils.State) {
	peersCount := uint32(len(state.Peers))
	totalMsgsCount := state.MyMsgCount

	for _, peer := range state.Peers {
		totalMsgsCount += peer.NumMsgs
	}

	state.MyDC = make([]uint64, totalMsgsCount)
	var i, j uint32

	// generates power sums of message_hashes
	// my_dc[i] := my_dc[i] (+) (my_msg_hashes[j] ** (i + 1))
	for j = 0; j < state.MyMsgCount; j++ {
		// generates 64 bit hash of my_message[j]
		state.MyMessagesHash[j] = shortHash(state.MyMessages[j])
		var pow uint64 = 1
		for i = 0; i < totalMsgsCount; i++ {
			var op1 = field.NewField(field.UInt64(state.MyDC[i]))
			pow = power(uint64(state.MyMessagesHash[j]), pow)
			var op2 = field.NewField(field.UInt64(pow))

			state.MyDC[i] = uint64(op1.Add(op2).Fp)
		}
	}

	// encode power sums
	// my_dc[i] := my_dc[i] (+) (sgn(my_id - p.id) (*) p.dicemix.get_field_element())
	for j = 0; j < peersCount; j++ {
		for i = 0; i < totalMsgsCount; i++ {
			var op1 = field.NewField(field.UInt64(state.MyDC[i]))
			var op2 = field.NewField(field.UInt64(state.Peers[j].Dicemix.GetFieldElement()))
			if state.MyID < state.Peers[j].ID {
				op2 = op2.Neg()
			}
			state.MyDC[i] = uint64(op1.Add(op2).Fp)
		}
	}

	log.Info("My DC-EXP vector = ", state.MyDC)
}

// Verify that every peer agrees to proceed
func (d *dcNet) VerifyProceed(state *utils.State) bool {
	var i uint32
	totalMsgsCount := state.MyMsgCount

	// if slot collision occured with one of my messages
	if !state.MyOk {
		return false
	}

	// if slot collision occured with one of my peers messages
	for _, peer := range state.Peers {
		totalMsgsCount += peer.NumMsgs
		if !peer.Ok {
			return false
		}
	}

	// if one of my peer provided wrong confirmation
	for i = 0; i < totalMsgsCount; i++ {
		s := shortHash(utils.BytesToBase58String(state.AllMessages[i]))
		if state.AllMsgHashes[i] != reduce(s) {
			return false
		}
	}
	return true
}

func shortHash(message string) uint64 {
	// NOTE: after DC-EXP roots would contain hash reduced into field
	// (as final result would be in field)
	return xhashes.FNV64(message)
}

// parameter sdhould be within uint64 range
func power(value, t uint64) uint64 {
	return uint64(field.NewField(field.UInt64(value)).Mul(field.NewField(field.UInt64(t))).Fp)
}

// reduces value into field range
func reduce(value uint64) uint64 {
	return uint64(field.NewField(field.UInt64(value)).Fp)
}
