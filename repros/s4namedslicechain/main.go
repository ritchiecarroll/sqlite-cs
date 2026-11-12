package main

import "fmt"

// sqlite-cs S4: github.com/remyoudompheng/bigfft `type nat []big.Word` and `type fermat nat` -- a defined type whose
// underlying type is another DEFINED slice type. Go gives fermat every slice operation.
type nat []uint64
type fermat nat

func main() {
	z := make(fermat, 3)
	z[0] = 7
	n := nat(z[1:])
	f := fermat(n)
	fmt.Println(len(z), z[0], len(n), len(f), cap(z[:2])) // Go: 3 7 2 2 3
}
