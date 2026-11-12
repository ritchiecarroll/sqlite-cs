package main

import (
	"fmt"

	sqlite3 "repro/s19b-namedstructrenamedimport/lib"
)

// sqlite-cs S19, narrowing (does NOT reproduce S19): busy.go's field store, through a pointer to a NAMED struct (not an
// alias of an anonymous struct, as Tsqlite3 is) from a package imported under another name.
func main() {
	db := new(sqlite3.Tdb)
	db.FbusyTimeout = 250
	fmt.Println(db.FbusyTimeout) // Go: 250
}
