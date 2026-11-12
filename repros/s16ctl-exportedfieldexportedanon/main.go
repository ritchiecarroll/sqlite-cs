package main

import (
	"fmt"

	"repro/s16ctl-exportedfieldexportedanon/lib"
)

// sqlite-cs S16 green control: the inner alias exported.
func main() {
	p := lib.Make(1, 9)
	fmt.Println(p.FfContext, p.Fu.FhInproc) // Go: 1 9
}
