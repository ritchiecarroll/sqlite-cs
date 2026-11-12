package main

import (
	"fmt"

	"repro/s16exportedfieldunexportedanon/lib"
)

// sqlite-cs S16: modernc.org/sqlite/lib TuserHPALETTE.Fu, an exported field typed with an unexported alias of an anonymous
// struct, read from another package.
func main() {
	p := lib.Make(1, 9)
	fmt.Println(p.FfContext, p.Fu.FhInproc) // Go: 1 9
}
