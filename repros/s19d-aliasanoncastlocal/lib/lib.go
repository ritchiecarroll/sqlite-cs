// Package lib holds the type as ccgo writes Tsqlite3: an alias of an anonymous struct, with an array field.
package lib

type Tdb = struct {
	FbusyTimeout int32
	FaDb         [2]int64
}
