package main

import (
	"fmt"
	"unsafe"

	sqlite3 "repro/s19d-aliasanoncastlocal/lib"
)

// sqlite-cs S19, narrowing (does NOT reproduce S19): busy.go's shape as far as it is known -- a local assigned only from
// `(*T)(unsafe.Pointer)`, T another package's alias of an anonymous struct with an array field, imported under another
// name, and a field store through the local. go2cs records it as `libꓸTdb` and builds; the driver's record reads
// `lib.Tsqlite3`, so an ingredient is still missing.
func busy(p unsafe.Pointer, t int32) {
	db := (*sqlite3.Tdb)(p)
	db.FbusyTimeout = t
}

func main() {
	buf := make([]int64, 3) // the struct's memory, never declared as the struct: busy.go only casts pointers to it
	busy(unsafe.Pointer(&buf[0]), 250)
	fmt.Println(int32(buf[0])) // Go: 250 (FbusyTimeout is at offset 0; little-endian)
}
