package dc

import (
	"../field"
	"../utils"
	"github.com/shomali11/util/xhashes"
)

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
