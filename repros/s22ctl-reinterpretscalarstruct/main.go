package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S22: modernc.org/sqlite/lib `var _sqlite3DigitPairs = *(*struct{ ... Fa [201]int8 ... })(unsafe.Pointer(&struct{...}{...}))`.
// This is the control for s22reinterpretarraystruct: the same reinterpretation with scalar fields.
var digits = *(*struct {
	Fa0, Fa1, Fa2, Fa3 int8
	Fpad               byte
})(unsafe.Pointer(&struct {
	f0, f1, f2, f3 int8
	_              byte
}{f0: '0', f1: '1', f2: '2', f3: '3'}))

func main() {
	fmt.Println(digits.Fa3) // Go: 51
}
