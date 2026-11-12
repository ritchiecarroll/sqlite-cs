package main

import (
	"fmt"
	"unsafe"

	"repro/s18b-foreignnamedarrayslice/mem"
)

// sqlite-cs S18, narrowing: the sqlite site's other ingredient, a named array type (libc's RawMem) declared in ANOTHER
// package, with an ASCII qualifier. It does NOT reproduce S18 (it builds, then fails as S5): S18 needs a non-ASCII
// character in the emitted type name (s18nonasciiptrcastslice), which sqlite gets from go2cs's renamed import Δlibc.
func main() {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	p := uintptr(unsafe.Pointer(&src[0]))
	v := make([]byte, 8)
	n := copy(v, (*mem.RawMem)(unsafe.Pointer(p))[:8:8])
	fmt.Println(n, v) // Go: 8 [1 2 3 4 5 6 7 8]
}
