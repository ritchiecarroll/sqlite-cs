package main

import "fmt"

// sqlite-cs S6 green control: the same conversion written explicitly as string(rune(b)).
type literalErr string

func (e literalErr) Error() string { return "strftime: unsupported literal: " + string(e) }

func check(b byte) error {
	if '0' <= b && b <= '9' {
		return literalErr(string(rune(b)))
	}
	return nil
}

func main() {
	fmt.Println(check('7'), check('x') == nil) // Go: strftime: unsupported literal: 7 true
}
