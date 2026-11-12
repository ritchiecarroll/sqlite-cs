package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S18: modernc.org/sqlite conn.go `copy(v, (*libc.RawMem)(unsafe.Pointer(p))[:len:len])`, emitted with go2cs's
// renamed import `Δlibc.RawMem`. Slicing a pointer cast whose C# type name holds a NON-ASCII character drops one character
// of the emitted cast (the converter measures the type name in bytes and cuts the expression in runes). Here the
// non-ASCII character is in the type's own name; s5ptrnamedarrayslice is the same code with an ASCII name.
type mémoire [1 << 20]byte

func main() {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	p := uintptr(unsafe.Pointer(&src[0]))
	v := make([]byte, 8)
	n := copy(v, (*mémoire)(unsafe.Pointer(p))[:8:8])
	fmt.Println(n, v) // Go: 8 [1 2 3 4 5 6 7 8]
}
