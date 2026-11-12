package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S17: modernc.org/sqlite/lib sqlite.go `(*(*TWalSehFindFrame)(__ccgo_ptr(bp))).Fpgno = pgno`: a field store
// through an EXPLICIT dereference of a pointer conversion. The control (s17ctl) writes `(*T)(p).F` (implicit dereference).
type frame struct {
	Pgno int32
	Read int32
}

func ptr(p uintptr) unsafe.Pointer { return unsafe.Pointer(p) }

func main() {
	f := new(frame)
	bp := uintptr(unsafe.Pointer(f))
	(*(*frame)(ptr(bp))).Pgno = 7
	(*(*frame)(ptr(bp))).Read += 3
	fmt.Println(f.Pgno, f.Read) // Go: 7 3
}
