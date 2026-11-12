package main

import "fmt"

// sqlite-cs S15: modernc.org/sqlite/lib `type AggInfo = TAggInfo` and `type TAggInfo = struct {...}`: an alias of an alias
// of an anonymous struct type, in ONE package (S2 is the cross-package twin).
type TAggInfo = struct {
	FnColumn int32
	FnFunc   int32
}

type AggInfo = TAggInfo

func main() {
	a := AggInfo{FnColumn: 2, FnFunc: 5}
	var t TAggInfo = a
	fmt.Println(a.FnColumn, t.FnFunc) // Go: 2 5
}
