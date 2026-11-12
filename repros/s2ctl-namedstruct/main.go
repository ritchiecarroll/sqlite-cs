package main

import (
	"fmt"
	"unsafe"

	"repro/s2ctl-namedstruct/lib"
)

// sqlite-cs S2 green control: the same alias, of a NAMED (defined) struct type in lib.
type cIndexConstraint = lib.Tindex_constraint

func main() {
	c := cIndexConstraint{FiColumn: 7, Fop: 2}
	fmt.Println(c.FiColumn, c.Fop, unsafe.Sizeof(c))
}
