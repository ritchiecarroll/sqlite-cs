package main

import "fmt"

// sqlite-cs S9 green control: the same aliases with names that are not C# keywords.
type (
	clong  = int32
	culong = uint32
)

func size(n int) clong { return clong(n) }

func main() {
	var u culong = 7
	fmt.Println(size(3), u, clong(-1)) // Go: 3 7 -1
}
