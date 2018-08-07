package dc

import (
	"fmt"

	"../field"
	"../utils"
	base58 "github.com/jbenet/go-base58"
	"github.com/shomali11/util/xhashes"
)

type dcNet struct {
	DC
}

// NewDCNetwork creates a new DC instance
func NewDCNetwork() DC {
	return &dcNet{}
}

// RunDCExponential -Runs DC-EXP
func (d *dcNet) RunDCSimple(state *utils.State) {
	// initaializing variables
	slots := make([]int, state.MyMsgCount)
	peersCount := uint32(len(state.Peers))
	myOK := true
	var i, j uint32

	for i := range slots {
		slots[i] = -1
	}

	for j = 0; j < state.MyMsgCount; j++ {
		index, count := -1, 0
		for i = 0; i < state.TotalMsgsCount; i++ {
			if state.AllMsgHashes[i] == reduce(state.MyMessagesHash[j]) {
				index, count = int(i), int(count+1)
			}
		}

		if count == 1 {
			slots[j] = index
		} else {
			myOK = false
		}
	}

	if !myOK {
		// Even though the run will be aborted (because we send my_ok = false), transmit the
		// message in a deterministic slot. This enables the peers to recompute our commitment.
		for i = 0; i < state.MyMsgCount; i++ {
			slots[i] = int(i)
		}
	}

	//  array of |totalMsgsCount| arrays of slot_size bytes, all initalized with 0
	state.DCSimpleVector = make([][]byte, state.TotalMsgsCount)

	for j = 0; j < state.TotalMsgsCount; j++ {
		state.DCSimpleVector[j] = make([]byte, 20)
	}

	for j = 0; j < state.MyMsgCount; j++ {
		fmt.Printf("\nBYTES = %v\n", base58StringToBytes(state.MyMessages[j]))
		state.DCSimpleVector[slots[j]] = base58StringToBytes(state.MyMessages[j])
	}

	fmt.Printf("\nSLOT's = %v\n", state.DCSimpleVector)

	for i = 0; i < peersCount; i++ {
		for j = 0; j < state.TotalMsgsCount; j++ {
			xorBytes(state.DCSimpleVector[j], state.DCSimpleVector[j], state.Peers[i].Dicemix.GetBytes(20))
		}
	}

	fmt.Printf("\nMY DC_SIMPLE[] = %v\n", state.DCSimpleVector)
}

func (d *dcNet) ResolveDCNet(state *utils.State) {
	var i, j uint32
	peersCount := uint32(len(state.Peers))
	state.AllMessages = state.DCSimpleVector

	for i = 0; i < peersCount; i++ {
		for j = 0; j < state.TotalMsgsCount; j++ {
			xorBytes(state.AllMessages[j], state.AllMessages[j], state.Peers[i].DCSimpleVector[j])
		}
	}

	fmt.Printf("\nMY RESOLVED DC NET VECTOR[] = %v\n", state.AllMessages)
}

// generates my_dc[]
func (d *dcNet) DeriveMyDCVector(state *utils.State) {
	peersCount := uint32(len(state.Peers))
	state.TotalMsgsCount = peersCount + state.MyMsgCount

	state.MyDC = make([]uint64, state.TotalMsgsCount)
	var i, j uint32

	for j = 0; j < state.MyMsgCount; j++ {
		state.MyMessagesHash[j] = shortHash(state.MyMessages[j])
		var pow uint64 = 1
		for i = 0; i < state.TotalMsgsCount; i++ {
			var op1 = field.NewField(field.UInt64(state.MyDC[i]))
			pow = power(uint64(state.MyMessagesHash[j]), pow)
			var op2 = field.NewField(field.UInt64(pow))

			state.MyDC[i] = uint64(op1.Add(op2).Fp)
		}
	}

	for j = 0; j < peersCount; j++ {
		for i = 0; i < state.TotalMsgsCount; i++ {
			var op1 = field.NewField(field.UInt64(state.MyDC[i]))
			var op2 = field.NewField(field.UInt64(state.Peers[j].Dicemix.GetFieldElement()))
			if state.MyID < state.Peers[j].ID {
				op2 = op2.Neg()
			}
			state.MyDC[i] = uint64(op1.Add(op2).Fp)
		}
	}

	fmt.Printf("\nMY DC_VECTOR[] = %v\n", state.MyDC)
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

func reduce(value uint64) uint64 {
	return uint64(field.NewField(field.UInt64(value)).Fp)
}

func bytesToBase58String(bytes []byte) string {
	return base58.Encode(bytes)
}

func base58StringToBytes(str string) []byte {
	return base58.Decode(str)
}
