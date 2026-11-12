package main

import (
	"fmt"
	"unsafe"
)

// sqlite-cs S14: modernc.org/libc qsort.go `*(*uint64)(unsafe.Pointer(p)), *(*uint64)(unsafe.Pointer(q)) =
// *(*uint64)(unsafe.Pointer(q)), *(*uint64)(unsafe.Pointer(p))`: a parallel assignment through dereferenced pointer casts.
func swap(p, q uintptr) {
	*(*uint64)(unsafe.Pointer(p)), *(*uint64)(unsafe.Pointer(q)) = *(*uint64)(unsafe.Pointer(q)), *(*uint64)(unsafe.Pointer(p))
}

func main() {
	a := []uint64{1, 2}
	swap(uintptr(unsafe.Pointer(&a[0])), uintptr(unsafe.Pointer(&a[1])))
	fmt.Println(a) // Go: [2 1]
}
