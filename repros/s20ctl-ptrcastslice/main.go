package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S20: modernc.org/memory memory.go:348 (UintptrCalloc) `b := ((*rawmem)(unsafe.Pointer(r)))[:size:size]`, with
// `type rawmem [1<<50 - 1]byte`. This is the control for s20parenptrcastslice: the same slice without the
// extra parentheses. The array here is unnamed, so S5 (a named array type) stays out of it.
func main() {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	p := uintptr(unsafe.Pointer(&src[0]))
	n := 8
	b := (*[1<<50 - 1]byte)(unsafe.Pointer(p))[:n:n]
	fmt.Println(len(b), b[0], b[7]) // Go: 8 1 8
}
