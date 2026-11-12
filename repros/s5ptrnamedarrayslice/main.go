package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S3/S5: modernc.org/memory's realloc, `copy((*rawmem)(unsafe.Pointer(r))[:size:size], (*rawmem)(unsafe.Pointer(p))[:size:size])`,
// with `type rawmem [1<<50 - 1]byte`: slicing a pointer to a NAMED array type views the memory as []byte, and the copy
// writes into r's memory (S3: a write through the view must reach the original; S5: the element type is byte).
type rawmem [1 << 20]byte

func main() {
	src := make([]byte, 16)
	dst := make([]byte, 16)
	for i := range src {
		src[i] = byte(i + 1)
	}
	p, r := uintptr(unsafe.Pointer(&src[0])), uintptr(unsafe.Pointer(&dst[0]))
	n := copy((*rawmem)(unsafe.Pointer(r))[:8:8], (*rawmem)(unsafe.Pointer(p))[:8:8])
	b := (*rawmem)(unsafe.Pointer(r))[:4]
	fmt.Println(n, dst[:9], len(b), b[3]) // Go: 8 [1 2 3 4 5 6 7 8 0] 4 4
}
