package main

import "fmt"

// sqlite-cs S17, narrowing: the same explicit-dereference field store `(*p).F = x`, but through a plain pointer
// variable, not a pointer conversion.
type frame struct {
	Pgno int32
	Read int32
}

func main() {
	p := new(frame)
	(*p).Pgno = 7
	(*p).Read += 3
	fmt.Println(p.Pgno, p.Read) // Go: 7 3
}
