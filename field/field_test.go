package field

import (
	"testing"
)

type testpair struct {
	data []UInt64
	res  UInt64
}

var additionTests = []testpair{
	{[]UInt64{7, 5}, 12},
	{[]UInt64{P - 2, 5}, 3},
	{[]UInt64{2193980333835211996, 621408416523297271}, 509545741144815316},
	{[]UInt64{18446744073709551615, 18446744073709551615}, 14},
	{[]UInt64{2305843009213693950, 3}, 2},
	{[]UInt64{2305843009213693950, 1}, 0},
}

var subtractionTests = []testpair{
	{[]UInt64{7, 5}, 2},
	{[]UInt64{4, 8}, P - 4},
	{[]UInt64{18446744073709551615, 18446744073709551615}, 0},
	{[]UInt64{P, 5}, P - 5},
}

var multiplicationTests = []testpair{
	{[]UInt64{4, 3}, 12},
	{[]UInt64{2239513929391938494, 1021644029483981869}, 619009326837417152},
	{[]UInt64{2305843009213693950, 5}, 2305843009213693946},
}

var negationTests = []testpair{
	{[]UInt64{4}, P - 4},
	{[]UInt64{P}, 0},
	{[]UInt64{P - 1}, 1},
	{[]UInt64{P + 5}, P - 5},
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
