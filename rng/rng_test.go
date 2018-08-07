package rng

import (
	"encoding/hex"
	"testing"
)

type testpair struct {
	seed string
	prg  string
}

var testcases = []testpair{
	{"0000000000000000000000000000000000000000000000000000000000000000", "76b8e0ada0f13d90"},
	{"0100000000000000000000000000000000000000000000000000000000000000", "c5d30a7ce1ec1193"},
}

func decodeString(key string) []byte {
	seed, _ := hex.DecodeString(key)
	return seed
}

func TestRng(t *testing.T) {
	for _, pair := range testcases {
		v := NewRng(decodeString(pair.seed))
		if hex.EncodeToString(v.chachaExpRng) != pair.prg {
			t.Error(
				"For", pair.seed,
				"expected", pair.prg[0],
				"got", v,
			)
		}
	}
}
