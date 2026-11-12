package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S20, explicit-dereference form: modernc.org/libc libc.go:197 (SetEnviron)
// `copy((*(*RawMem)(unsafe.Pointer(s)))[:len(v):len(v)], v)`: a pointer-to-array cast DEREFERENCED, then sliced. The control (s20ctl) is the
// same slice taken straight from the pointer. The array here is unnamed, so S5 (a named array type) stays out of it.
func main() {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	p := uintptr(unsafe.Pointer(&src[0]))
	n := 8
	b := (*(*[1<<50 - 1]byte)(unsafe.Pointer(p)))[:n:n]
	fmt.Println(len(b), b[0], b[7]) // Go: 8 1 8
}
