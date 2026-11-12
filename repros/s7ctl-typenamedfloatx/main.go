package main

import "fmt"

// sqlite-cs S7 green control: the same type and literal, named floatx.
// named floatx, built with a composite literal.
type floatx struct {
	n, fracBits int
}

func (f *floatx) normalize() { f.fracBits++ }

func newFloat(n, fracBits int) floatx {
	f := floatx{n: n, fracBits: fracBits}
	f.normalize()
	return f
}

func main() {
	f := newFloat(3, 4)
	fmt.Println(f.n, f.fracBits) // Go: 3 5
}
