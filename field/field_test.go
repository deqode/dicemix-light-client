package field

import (
	"testing"
)

type testpair struct {
	data []uint64
	res  uint64
}

var additionTests = []testpair{
	{[]uint64{7, 5}, 12},
	{[]uint64{P.Value() - 2, 5}, 3},
	{[]uint64{2193980333835211996, 621408416523297271}, 509545741144815316},
	{[]uint64{18446744073709551615, 18446744073709551615}, 14},
	{[]uint64{2305843009213693950, 3}, 2},
	{[]uint64{2305843009213693950, 1}, 0},
}

var subtractionTests = []testpair{
	{[]uint64{7, 5}, 2},
	{[]uint64{4, 8}, P.Value() - 4},
	{[]uint64{18446744073709551615, 18446744073709551615}, 0},
	{[]uint64{P.Value(), 5}, P.Value() - 5},
}

var multiplicationTests = []testpair{
	{[]uint64{4, 3}, 12},
	{[]uint64{2239513929391938494, 1021644029483981869}, 619009326837417152},
	{[]uint64{2305843009213693950, 5}, 2305843009213693946},
}

var negationTests = []testpair{
	{[]uint64{4}, P.Value() - 4},
	{[]uint64{P.Value()}, 0},
	{[]uint64{P.Value() - 1}, 1},
	{[]uint64{P.Value() + 5}, P.Value() - 5},
}

func TestAddition(t *testing.T) {
	for _, pair := range additionTests {
		v := NewField(pair.data[0]).Add(NewField(pair.data[1]))
		if v != NewField(pair.res) {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", v,
			)
		}
	}
}

func TestSubtraction(t *testing.T) {
	for _, pair := range subtractionTests {
		v := NewField(pair.data[0]).Sub(NewField(pair.data[1]))
		if v != NewField(pair.res) {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", v,
			)
		}
	}
}

func TestMultiplication(t *testing.T) {
	for _, pair := range multiplicationTests {
		v := NewField(pair.data[0]).Mul(NewField(pair.data[1]))
		if v != NewField(pair.res) {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", v,
			)
		}
	}
}

func TestNegation(t *testing.T) {
	for _, pair := range negationTests {
		v := NewField(pair.data[0]).Neg()
		if v != NewField(pair.res) {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", v,
			)
		}
	}
}
