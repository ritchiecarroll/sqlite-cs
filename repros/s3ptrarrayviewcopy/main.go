package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S3: a write through a slice of a pointer-to-array cast must reach the memory it views (realloc, memory.go:502).

func main() {
	src := make([]byte, 16)
	dst := make([]byte, 16)
	for i := range src {
		src[i] = byte(i + 1)
	}
	p, r := uintptr(unsafe.Pointer(&src[0])), uintptr(unsafe.Pointer(&dst[0]))
	n := copy((*[1 << 20]byte)(unsafe.Pointer(r))[:8:8], (*[1 << 20]byte)(unsafe.Pointer(p))[:8:8])
	b := (*[1 << 20]byte)(unsafe.Pointer(r))[:4]
	fmt.Println(n, dst[:9], len(b), b[3]) // Go: 8 [1 2 3 4 5 6 7 8 0] 4 4
}
