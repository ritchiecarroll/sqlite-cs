package main

import (
	"fmt"
	"unsafe"

	"repro/s12ctl-foreignnamed/stdio"
)

// sqlite-cs S12 green control: the same array typed with the named struct itself, not its alias.
// `X_iob [3]stdio.X_iobuf`: another package's ALIAS named by qualifier, here in an element-address position.
var X_iob [3]stdio.X_iobuf

func main() {
	m := map[uintptr]int32{}
	for i := range X_iob {
		m[uintptr(unsafe.Pointer(&X_iob[i]))] = int32(i)
	}
	X_iob[1].Fcnt = 4
	fmt.Println(len(m), X_iob[1].Fcnt, unsafe.Sizeof(X_iob[0])) // Go: 3 4 16
}
