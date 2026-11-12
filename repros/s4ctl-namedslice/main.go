package main

import "fmt"

// sqlite-cs S4 green control: fermat declared directly over the slice type.
type nat []uint64
type fermat []uint64

func main() {
	z := make(fermat, 3)
	z[0] = 7
	n := nat(z[1:])
	f := fermat(n)
	fmt.Println(len(z), z[0], len(n), len(f), cap(z[:2])) // Go: 3 7 2 2 3
}
