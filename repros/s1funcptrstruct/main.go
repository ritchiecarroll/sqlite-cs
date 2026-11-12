package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S1: modernc.org/sqlite/lib's function-pointer call, `(*(*func(...))(unsafe.Pointer(&struct{ uintptr }{fp})))(...)`.
// ccgo keeps a C function pointer as a uintptr (__ccgo_fp reads an interface's data word), and a call rebuilds a Go func
// value from that word through an anonymous struct{ uintptr }.
func fpOf(f interface{}) uintptr {
	type iface [2]uintptr
	return (*iface)(unsafe.Pointer(&f))[1]
}

func twice(n int) { fmt.Println("twice", 2*n) }

func main() {
	fp := fpOf(twice)
	if fp != 0 {
		(*(*func(int))(unsafe.Pointer(&struct{ uintptr }{fp})))(21)
	}
}
