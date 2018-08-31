package ecdh

import (
	"bytes"
	"testing"
)

type sharedSecretTestPair struct {
	keys   [][]byte
	secret []byte
}

var sharedSecretTests = []sharedSecretTestPair{
	{
		[][]byte{
			{176, 171, 9, 141, 190, 148, 191, 54, 202, 77, 181, 13, 110, 89, 170, 164, 100, 150, 213, 73, 206, 86, 53, 119, 249, 126, 40, 254, 107, 32, 7, 125},
			{166, 96, 104, 213, 152, 216, 8, 63, 162, 73, 88, 69, 30, 23, 148, 213, 202, 181, 121, 219, 228, 238, 1, 249, 244, 72, 184, 97, 40, 241, 22, 106},
		}, []byte{99, 223, 247, 56, 31, 76, 52, 54, 68, 53, 162, 121, 2, 83, 47, 21, 132, 250, 4, 253, 6, 143, 186, 245, 7, 243, 247, 169, 88, 209, 239, 30},
	},

	{
		[][]byte{
			{128, 81, 6, 59, 12, 192, 217, 123, 71, 114, 134, 244, 159, 118, 237, 42, 230, 155, 72, 185, 237, 21, 93, 174, 105, 170, 213, 155, 23, 231, 118, 123},
			{165, 236, 127, 87, 7, 227, 231, 102, 47, 253, 108, 228, 222, 223, 147, 102, 184, 209, 227, 64, 67, 21, 204, 1, 254, 60, 187, 183, 209, 125, 223, 41},
		}, []byte{99, 223, 247, 56, 31, 76, 52, 54, 68, 53, 162, 121, 2, 83, 47, 21, 132, 250, 4, 253, 6, 143, 186, 245, 7, 243, 247, 169, 88, 209, 239, 30},
	},
}

func TestSharedSecret(t *testing.T) {
	ecdh := NewCurve25519ECDH()
	for _, pair := range sharedSecretTests {
		privateKey, _ := ecdh.UnmarshalSK(pair.keys[0])
		publicKey, _ := ecdh.Unmarshal(pair.keys[1])
		secret, _ := ecdh.GenerateSharedSecret(privateKey, publicKey)

		if !bytes.Equal(pair.secret, secret) {
			t.Error(
				"For", pair.keys,
				"expected", pair.secret,
				"got", secret,
			)
		}
	}
}
