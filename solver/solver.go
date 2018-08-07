package solver

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lflint -lgmp
// #include <stdlib.h>
// #include "solver_flint.h"
import "C"
import (
	"fmt"
	"strconv"
	"unsafe"

	"../field"
)

// Solve -- solves the generated DC-COMBINED[] to obtain MESSAGES HASHES (if no error exists)
// else return NULL.
// Runs on server side, solves polynomial
// and returns generated roots to all clients for verification
func Solve(dcCombined []uint64, count int) []uint64 {

	fmt.Printf("%v %v\n", dcCombined, count)

	// would contain results generated through solver_flint
	outMessages := make([]string, count)

	// contains dc-combined[]
	sums := make([]string, count)

	// value of our Field range (P)
	prime := C.CString(fmt.Sprint(field.P))

	for i := 0; i < count; i++ {
		sums[i] = fmt.Sprint(dcCombined[i])
	}

	defer C.free(unsafe.Pointer(prime))

	// Allocate memory to outMessages[] and convert it into C String
	argcOutMessages := C.int(count)
	valueOutMessages := (*[0xfff]*C.char)(C.allocArgv(argcOutMessages))
	defer C.free(unsafe.Pointer(valueOutMessages))

	for i, arg := range outMessages {
		valueOutMessages[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(valueOutMessages[i]))
	}

	// Allocate memory to sums[] and convert it into C String
	argsSums := C.int(count)
	valueSums := (*[0xfff]*C.char)(C.allocArgv(argsSums))
	defer C.free(unsafe.Pointer(valueSums))

	for i, arg := range sums {
		valueSums[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(valueSums[i]))
	}

	// Call solve() of solver_flint.cpp
	// returns string[] containing sorted Messge Hashes (if successful)
	// else returns NULL
	var x **C.char = C.solve(C.int(count), (**C.char)(unsafe.Pointer(valueOutMessages)), prime, (**C.char)(unsafe.Pointer(valueSums)))

	// Convert C-Style string[] to Go-Style uint64[]
	var messages = goStrings(C.int(count), x)

	return messages
}

// Convert C-Style string[] to Go-Style uint64[]
func goStrings(argc C.int, argv **C.char) []uint64 {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]uint64, length)
	for i, s := range tmpslice {
		gostrings[i], _ = strconv.ParseUint(C.GoString(s), 10, 64)
	}
	return gostrings
}
