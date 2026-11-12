package main

import (
	"fmt"

	"repro/s19c-childpkgrenamedimport/drv"
)

// sqlite-cs S19, narrowing (does NOT reproduce S19): the field store of s19b-namedstructrenamedimport, but in a library package that imports its own
// CHILD package under another name (modernc.org/sqlite -> sqlite3 "modernc.org/sqlite/lib").
func main() {
	fmt.Println(drv.Busy(250)) // Go: 250
}
