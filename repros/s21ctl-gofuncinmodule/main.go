package main

import "fmt"

// sqlite-cs S21: golang.org/x/sys/cpu declares `func cpuid(eaxArg, ecxArg uint32) (eax, ebx, ecx, edx uint32)` in
// cpu_gc_x86.go and implements it in cpu_gc_x86.s; x/sys/cpu's init calls it, so every program importing modernc.org/libc
// runs it at start-up. This is the control for s21asmfuncinmodule: the same function written in Go.
func add(a, b uint32) uint32 { return a + b }

func main() {
	fmt.Println(add(2, 3)) // Go: 5
}
