package main

import "fmt"

// sqlite-cs S7: modernc.org/mathutil binarylog.go `type float struct{...}` and `f := float{n: n, ...}`: a user struct type
// named float, built with a composite literal.
type float struct {
	n, fracBits int
}

func (f *float) normalize() { f.fracBits++ }

func newFloat(n, fracBits int) float {
	f := float{n: n, fracBits: fracBits}
	f.normalize()
	return f
}

func main() {
	f := newFloat(3, 4)
	fmt.Println(f.n, f.fracBits) // Go: 3 5
}
