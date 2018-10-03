package dc

import (
	"github.com/techracers-blockchain/dicemix-light-client/field"
	"github.com/techracers-blockchain/dicemix-light-client/utils"

	"github.com/shomali11/util/xhashes"
)

func obtainSlots(state *utils.State, totalMsgsCount uint32) ([]int, bool) {
	var i, j uint32
	slots := make([]int, state.MyMsgCount)
	ok := true

	// initialize slots
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
			ok = false
		}
	}
	return slots, ok
}

func shortHash(message string) uint64 {
	// NOTE: after DC-EXP roots would contain hash reduced into field
	// (as final result would be in field)
	return xhashes.FNV64(message)
}

// parameter sdhould be within uint64 range
func power(value, t uint64) uint64 {
	return field.NewField(value).Mul(field.NewField(t)).Value()
}

// reduces value into field range
func reduce(value uint64) uint64 {
	return field.NewField(value).Value()
}

// returns total numbers of messages
// my-msg-count +  ∑(peers.msg-count)
func messageCount(count uint32, peers []utils.Peers) uint32 {
	// calculate total number of messages
	for _, peer := range peers {
		count += peer.NumMsgs
	}
	return count
}
