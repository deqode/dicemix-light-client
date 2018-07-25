package field

import (
	"log"
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

func asLimbs(x UInt64) (uint32, uint32) {
	return uint32(x >> 32), uint32(x)
}

func NewField(value UInt64) Field {
	return Field{value.reduceOnce().reduceOnceAssert()}
}

func (self Field) Neg() Field {
	return Field{(P - self.fp).reduceOnce().reduceOnceAssert()}
}

func (self Field) Add(op2 Field) Field {
	return Field{(self.fp + op2.fp).reduceOnce().reduceOnceAssert()}
}

func (self *Field) AddAssign(op2 Field) {
	*self = self.Add(op2)
}

func (self Field) Sub(op2 Field) Field {
	if op2.fp > self.fp {
		return Field{(P - op2.fp + self.fp).reduceOnce().reduceOnceAssert()}
	}
	return Field{(self.fp - op2.fp).reduceOnce().reduceOnceAssert()}
}

func (self *Field) SubAssign(op2 Field) {
	*self = self.Sub(op2)
}
