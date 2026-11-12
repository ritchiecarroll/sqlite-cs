package main

import "fmt"

// sqlite-cs S13 green control: the same terminating switch with no fallthrough (default returns).
func lock(typ int) int32 {
	switch typ {
	default:
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
