package main

import "fmt"

// sqlite-cs S8: modernc.org/mathutil int.go `func (x *Int128) SetBigInt(y *big.Int) (r Int128, err error)` uses the
// RECEIVER POINTER as a value (`fmt.Errorf("%T.SetInt: overflow", x)`), and `NewInt128FromFloat32` calls r.SetBigInt(bi)
// on its addressable named result r (Go takes &r implicitly).
type Int128 struct{ Lo, Hi int64 }

func (x *Int128) SetBig(v int64) (r Int128, err error) {
	if v > 1000 {
		return *x, fmt.Errorf("%T.SetBig: overflow", x)
	}
	x.Lo, x.Hi = v, v>>63
	return *x, nil
}

func fromInt(v int64) (r Int128) {
	r.SetBig(v)
	return r
}

func main() {
	_, err := new(Int128).SetBig(5000)
	fmt.Println(fromInt(-5), fromInt(9), err) // Go: {-5 -1} {9 0} *main.Int128.SetBig: overflow
}
