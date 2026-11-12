package main

import "fmt"

// sqlite-cs S15 green control: the single alias of the anonymous struct, with no second alias.
type TAggInfo = struct {
	FnColumn int32
	FnFunc   int32
}

func main() {
	a := TAggInfo{FnColumn: 2, FnFunc: 5}
	var t TAggInfo = a
	fmt.Println(a.FnColumn, t.FnFunc) // Go: 2 5
}
