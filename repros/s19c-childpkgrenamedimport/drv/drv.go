// Package drv imports its own child package under another name, as modernc.org/sqlite imports
// `sqlite3 "modernc.org/sqlite/lib"`.
package drv

import sqlite3 "repro/s19c-childpkgrenamedimport/drv/lib"

// Busy stores a field through a pointer to the child package's struct (busy.go's `db.FbusyTimeout = tmout`).
func Busy(t int32) int32 {
	db := new(sqlite3.Tdb)
	db.FbusyTimeout = t
	return db.FbusyTimeout
}
