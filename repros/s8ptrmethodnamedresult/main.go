package main

import "fmt"

// sqlite-cs S8: modernc.org/mathutil int.go `func NewInt128FromFloat32(n float32) (r Int128) { ...; r.SetBigInt(bi); ... }`:
// a pointer-receiver method called on an addressable NAMED RESULT of value type (Go takes &r implicitly).
type Int128 struct{ Lo, Hi int64 }

func (x *Int128) SetBig(v int64) (r Int128, err error) {
	x.Lo, x.Hi = v, v>>63
	return *x, nil
}

func fromInt(v int64) (r Int128) {
	r.SetBig(v)
	return r
}

func main() {
	fmt.Println(fromInt(-5), fromInt(9)) // Go: {-5 -1} {9 0}
}
