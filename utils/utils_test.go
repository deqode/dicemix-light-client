package utils

import (
	"bytes"
	"testing"
)

type testpair struct {
	bytes []byte
	str   string
}

var tests = []testpair{
	{
		[]byte{191, 196, 191, 88, 228, 52, 126, 202, 91, 18, 50, 157, 255, 25, 125, 91, 8, 75, 201, 85},
		"3fxTZkRBKNj3KmJVkM3du4BP15e4",
	},
	{
		[]byte{47, 65, 205, 229, 8, 1, 35, 24, 239, 112, 76, 187, 142, 253, 154, 143, 217, 141, 150, 142},
		"fBkugJ6GtKuA5ez1Xt29hxfMjfj",
	},
	{
		[]byte{87, 228, 213, 202, 230, 147, 141, 171, 155, 118, 147, 79, 17, 202, 203, 130, 238, 42, 42, 119},
		"2E2FWkeU8Ku4trPYfrCPLxaU79DY",
	},
}

func TestBytesToBase58String(t *testing.T) {
	for _, pair := range tests {
		output := BytesToBase58String(pair.bytes)

		if output != pair.str {
			t.Error(
				"For", pair.bytes,
				"expected", pair.str,
				"got", output,
			)
		}
	}
}

func TestBase58StringToBytes(t *testing.T) {
	for _, pair := range tests {
		output := Base58StringToBytes(pair.str)

		if !bytes.Equal(output, pair.bytes) {
			t.Error(
				"For", pair.str,
				"expected", pair.bytes,
				"got", output,
			)
		}
	}
}
