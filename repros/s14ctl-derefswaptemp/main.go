package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S14 green control: the same swap written through a temporary.
func swap(p, q uintptr) {
	t := *(*uint64)(unsafe.Pointer(p))
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(q))
	*(*uint64)(unsafe.Pointer(q)) = t
}

func main() {
	a := []uint64{1, 2}
	swap(uintptr(unsafe.Pointer(&a[0])), uintptr(unsafe.Pointer(&a[1])))
	fmt.Println(a) // Go: [2 1]
}
