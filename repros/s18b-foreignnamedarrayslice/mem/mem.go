// Package mem holds the named array type, as modernc.org/libc holds RawMem.
package mem

type RawMem [1 << 20]byte
