package main

import (
	"fmt"
	"time"
)

// sqlite-cs S10 green control: the same type through the plain import "time".
type tzLocation struct {
	tz  string
	loc *time.Location
}

func main() {
	z := tzLocation{tz: "UTC", loc: time.UTC}
	fmt.Println(z.tz, z.loc.String()) // Go: UTC UTC
}
