package main

import "fmt"

// sqlite-cs S8 (variant): modernc.org/mathutil int.go has SetBigInt on BOTH *Int128 and *Uint128, each returning named
// results, and calls r.SetBigInt(bi) on an addressable named result of each type.
type Int128 struct{ Lo, Hi int64 }
type Uint128 struct{ Lo, Hi uint64 }

func (x *Int128) SetBig(v int64) (r Int128, err error)    { x.Lo, x.Hi = v, v>>63; return *x, nil }
func (x *Uint128) SetBig(v uint64) (r Uint128, err error) { x.Lo = v; return *x, nil }

func fromInt(v int64) (r Int128) {
	if v != 0 {
		r.SetBig(v)
	}
	return r
}

func fromUint(v uint64) (r Uint128) {
	r.SetBig(v)
	return r
}

func main() {
	fmt.Println(fromInt(-5), fromUint(9)) // Go: {-5 -1} {9 0}
}
