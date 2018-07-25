package field

import (
	"log"

	"github.com/cznic/mathutil"
)

type UInt64 uint64

const P UInt64 = (1 << 61) - 1

type Field struct {
	fp UInt64
}

func (self UInt64) reduceOnce() UInt64 {
	var value = (self & P) + (self >> 61)
	if value == P {
		return 0
	}
	return value
}

func (self UInt64) reduceOnceAssert() UInt64 {
	var res UInt64 = self.reduceOnce()
	if res >= P {
		log.Fatal("Error: Expected result should be less than field size")
	}
	return res
}

func (self UInt64) reduceOnceMul(op2 UInt64) UInt64 {
	var value UInt64 = (self << 3) | (op2 >> 61)
	value = (op2 & P) + value
	if value == P {
		return 0
	}
	return value
}

func asLimbs(x UInt64) (uint32, uint32) {
	return uint32(x >> 32), uint32(x)
}

// NewField reduces initial value in the field
func NewField(value UInt64) Field {
	return Field{value.reduceOnce().reduceOnceAssert()}
}

// Neg negates number in the field
func (self Field) Neg() Field {
	return Field{(P - self.fp).reduceOnce().reduceOnceAssert()}
}

// Add adds two field elements and reduces the resulting number in the field
func (self Field) Add(op2 Field) Field {
	return Field{(self.fp + op2.fp).reduceOnce().reduceOnceAssert()}
}

// AddAssign works same as Add, assigns final value to self
func (self *Field) AddAssign(op2 Field) {
	*self = self.Add(op2)
}

// Sub subtracts two field elements and reduces the resulting number in the field
func (self Field) Sub(op2 Field) Field {
	if op2.fp > self.fp {
		return Field{(P - op2.fp + self.fp).reduceOnce().reduceOnceAssert()}
	}
	return Field{(self.fp - op2.fp).reduceOnce().reduceOnceAssert()}
}

// SubAssign works same as Sub, assigns final value to self
func (self *Field) SubAssign(op2 Field) {
	*self = self.Sub(op2)
}

// Mul muliplies two field elements and reduces the resulting number in the field
func (self Field) Mul(op2 Field) Field {
	var high, low uint64 = mathutil.MulUint128_64(uint64(self.fp), uint64(op2.fp))
	var rh, rl = UInt64(high), UInt64(low)
	var res = rh.reduceOnceMul(rl).reduceOnceAssert()
	return Field{res}
}

// MulAssign works same as Mul, assigns final value to self
func (self *Field) MulAssign(op2 Field) {
	*self = self.Mul(op2)
}
