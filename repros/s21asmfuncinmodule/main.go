package main

import "fmt"

// sqlite-cs S21: golang.org/x/sys/cpu declares `func cpuid(eaxArg, ecxArg uint32) (eax, ebx, ecx, edx uint32)` in
// cpu_gc_x86.go and implements it in cpu_gc_x86.s; x/sys/cpu's init calls it, so every program importing modernc.org/libc
// runs it at start-up. A function implemented in assembly in a module outside the Go standard library. The control
// (s21ctl) has the same function written in Go.
func add(a, b uint32) uint32

func main() {
	fmt.Println(add(2, 3)) // Go: 5
}
