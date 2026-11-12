package main

import (
	"fmt"
	gotime "time"
)

// sqlite-cs S10: modernc.org/libc libc.go `gotime "time"` and `loc *gotime.Location`: a std type that go2cs Δ-renames
// (time.Location), reached through a RENAMED import.
type tzLocation struct {
	tz  string
	loc *gotime.Location
}

func main() {
	z := tzLocation{tz: "UTC", loc: gotime.UTC}
	fmt.Println(z.tz, z.loc.String()) // Go: UTC UTC
}
