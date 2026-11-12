package main

import "fmt"

// sqlite-cs S13: modernc.org/libc pthread.go `(*mutex).lock`: a function ending in a TERMINATING switch whose default
// falls through and whose clauses return or loop forever, so Go needs no final return.
func lock(typ int) int32 {
	switch typ {
	default:
		fallthrough
	case 0:
		return 10
	case 1:
		for n := 0; ; n++ {
			if n == 3 {
				return int32(20 + n)
			}
		}
	}
}

func main() {
	fmt.Println(lock(0), lock(1), lock(7)) // Go: 10 23 10
}
