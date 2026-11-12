package main

import "fmt"

// sqlite-cs S6: github.com/ncruces/go-strftime `return literalErr(b)`, b a byte, literalErr a named string type: Go's
// string(rune(b)) conversion into a named string type.
type literalErr string

func (e literalErr) Error() string { return "strftime: unsupported literal: " + string(e) }

func check(b byte) error {
	if '0' <= b && b <= '9' {
		return literalErr(b)
	}
	return nil
}

func main() {
	fmt.Println(check('7'), check('x') == nil) // Go: strftime: unsupported literal: 7 true
}
