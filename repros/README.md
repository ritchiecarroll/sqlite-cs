# Reproductions

Each failure class found while converting modernc.org/sqlite is reduced to a one-file Go module here, with a green
control that differs on one axis. Every module runs natively (`go run .` prints the line in its comment). Each class is a
gap in go2cs, and the fixes go there; this table says what each one does at the go2cs version stated.

Results are at go2cs master `4490a52244`, measured on windows/amd64 by converting each module with `go2cs -recurse`,
building the C#, running it and comparing its output with Go's.

| Row | Class | Repro (red) | At 4490a52244 | Control (green) | Where it appears in sqlite |
|---|---|---|---|---|---|
| S1 | A C function pointer kept as a `uintptr` and called by rebuilding a Go func value from `struct{ uintptr }`; the converter spells the reinterpret's source type as raw Go | `s1funcptrstruct` | conversion refused: `struct{uintptr}` emitted as raw Go | | 503 sites (lib, libc). Behind it, a delegate cannot be rebuilt from a number: a runtime design gap |
| S2 | An alias in another package of an alias of an anonymous struct | `s2aliasanonstruct` | conversion refused: the struct emitted as raw Go | `s2ctl-namedstruct` | 4 sites (driver's vtab.go, pagecache.go) |
| S3 | A slice of a pointer-to-array cast COPIES the memory it views, so a write through it is lost | `s3ptrarrayviewcopy` | runs, wrong output: `8 [0 0 0 0 0 0 0 0 0] 4 0` (Go `8 [1 2 3 4 5 6 7 8 0] 4 4`) | | 20 sites, the converter warns at each; the allocator's realloc (memory.go:502) among them |
| S4 | A defined type over another defined slice type (`type fermat nat`, `type nat []big.Word`) loses slice operations and its conversion to `nat` | `s4namedslicechain` | build fails (CS1503, CS0021, ...) | `s4ctl-namedslice` | bigfft, 131 errors |
| S5 | Slicing a pointer to a NAMED array type (`type rawmem [N]byte`) gives a slice of the array type, not of `byte` | `s5ptrnamedarrayslice` | runs, then `ArgumentException: Cannot use type 'rawmem'` | | memory (5 sites), libc's RawMem (about 9 sites); the same sites as S3 |
| S6 | `literalErr(b)`, a byte converted to a named string type (Go's `string(rune(b))`) | `s6bytetonamedstring` | build fails (CS0030) | `s6ctl-stringtonamedstring` | go-strftime, 1 site |
| S7 | A user struct type named `float`: its composite literal is emitted as `float64` | `s7typenamedfloat` | build fails (CS1739, CS1929) | `s7ctl-typenamedfloatx` | mathutil, 1 site |
| S8 | A pointer method that uses its receiver as a value (`fmt.Errorf("%T", x)`), called on an addressable named result: the address is not taken | `s8c-recvescapesnamedresult` | build fails (CS1929) | `s8ptrmethodnamedresult` (receiver not used as a value) | mathutil, 4 sites |
| S9 | A type alias whose name is a C# keyword (`long`, `ulong`, `byte`) | `s9aliaskeywordname` | build fails (parse errors) | `s9ctl-aliasplainname` | libc (`long`, `ulong`), lib (`byte`) |
| S10 | A std identifier go2cs renames (`time.Location`, `Local`, `UTC`, `Second`), reached through a renamed import (`gotime "time"`) | `s10aliasedimportrenamed` | build fails (CS0426) | `s10ctl-plainimport` | libc, 9 sites |
| S12 | Another package's alias named by qualifier (`stdio.FILE`, `type FILE = X_iobuf`) | `s12foreignaliasqualified` | build fails (CS0426) | `s12ctl-foreignnamed` | libc, 1 site |
| S13 | A function ending in a terminating switch whose default falls through: C# sees a reachable end | `s13terminatingfallthrough` | build fails (CS0161) | `s13ctl-terminatingswitch` | libc pthread, 3 functions |
| S14 | A parallel assignment through dereferenced pointer casts (a swap) | `s14derefswap` | build fails (CS0131) | `s14ctl-derefswaptemp` | libc qsort, 3 swaps |
| S15 | An alias of an alias of an anonymous struct in one package, the alias declared before and in another file from its target | `s15b-aliasforwardotherfile` | build fails (a `/` in the alias target) | `s15ctl-aliasanonstruct` | lib, 3,330 sites |
| S16 | An exported field whose type is an unexported alias of an anonymous struct, reached only as that field | `s16b-fieldonlyuse` | build fails (CS0050) | `s16ctl-exportedfieldexportedanon` | lib, 8 types and the types they reach |

Rows S2, S9, S12, S15 and S16 are likely one family (alias emission); S5 shares S3's sites. Variants that did NOT reproduce
their class are kept as negative evidence: `s8b-twotypessamemethod`, `s15aliaschainanonstruct`,
`s16exportedfieldunexportedanon`.
