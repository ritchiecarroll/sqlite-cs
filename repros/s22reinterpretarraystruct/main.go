package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S22: modernc.org/sqlite/lib `var _sqlite3DigitPairs = *(*struct{ FforceAlignment [0]int16; Fa [201]int8;
// F__ccgo_pad2 [1]byte })(unsafe.Pointer(&struct{ f [201]int8; _ [1]byte }{f: [201]int8{...}}))`, ccgo's alignment
// wrapper: one struct read as another through unsafe.Pointer, where both hold ARRAYS. The control (s22ctl) is the same
// reinterpretation with scalar fields.
var digits = *(*struct {
	FforceAlignment [0]int16
	Fa              [4]int8
	Fpad            [1]byte
})(unsafe.Pointer(&struct {
	f [4]int8
	_ [1]byte
}{f: [4]int8{'0', '1', '2', '3'}}))

func main() {
	fmt.Println(digits.Fa[3]) // Go: 51
}
