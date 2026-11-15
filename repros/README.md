# Reproductions

Each failure class found while converting modernc.org/sqlite is reduced to a small Go module here (one file, or two
packages where the class needs a second one), with a green
control that differs on one axis. Every module runs natively (`go run .` prints the line in its comment). Each class is a
gap in go2cs, and the fixes go there; this table says what each one does at the go2cs version stated.

Results are at go2cs master `4490a52244`, measured on windows/amd64 by converting each module with `go2cs -recurse`,
building the C#, running it and comparing its output with Go's. The *Fixed in go2cs* column names the go2cs commit that
fixes the class, where one has landed on go2cs master; those fixes are not yet re-measured here, so the *At* column
still reads as measured.

| Row | Class | Repro (red) | At 4490a52244 | Control (green) | Fixed in go2cs | Where it appears in sqlite |
|---|---|---|---|---|---|---|
| S1 | A C function pointer kept as a `uintptr` and called by rebuilding a Go func value from `struct{ uintptr }`; the converter spells the reinterpret's source type as raw Go | `s1funcptrstruct` | conversion refused: `struct{uintptr}` emitted as raw Go | | `aef262d42f` | 503 sites (lib, libc). Behind it, a delegate cannot be rebuilt from a number: a runtime design gap |
| S2 | An alias in another package of an alias of an anonymous struct | `s2aliasanonstruct` | conversion refused: the struct emitted as raw Go | `s2ctl-namedstruct` | `cf8f2cc8bb` | 4 sites (driver's vtab.go, pagecache.go) |
| S3 | A slice of a pointer-to-array cast COPIES the memory it views, so a write through it is lost | `s3ptrarrayviewcopy` | runs, wrong output: `8 [0 0 0 0 0 0 0 0 0] 4 0` (Go `8 [1 2 3 4 5 6 7 8 0] 4 4`) | | `996ce6c9bf` | 20 sites, the converter warns at each; the allocator's realloc (memory.go:502) among them |
| S4 | A defined type over another defined slice type (`type fermat nat`, `type nat []big.Word`) loses slice operations and its conversion to `nat` | `s4namedslicechain` | build fails (CS1503, CS0021, ...) | `s4ctl-namedslice` | `c3bda2e3b5` | bigfft, 131 errors |
| S5 | Slicing a pointer to a NAMED array type (`type rawmem [N]byte`) gives a slice of the array type, not of `byte` | `s5ptrnamedarrayslice` | runs, then `ArgumentException: Cannot use type 'rawmem'` | | `b72be15017` (the element type) | memory (5 sites), libc's RawMem (about 9 sites); the same sites as S3 |
| S6 | `literalErr(b)`, a byte converted to a named string type (Go's `string(rune(b))`) | `s6bytetonamedstring` | build fails (CS0030) | `s6ctl-stringtonamedstring` | `0ed7e0fc09` | go-strftime, 1 site |
| S7 | A user struct type named `float`: its composite literal is emitted as `float64` | `s7typenamedfloat` | build fails (CS1739, CS1929) | `s7ctl-typenamedfloatx` | `c3bda2e3b5` | mathutil, 1 site |
| S8 | A pointer method that uses its receiver as a value (`fmt.Errorf("%T", x)`), called on an addressable named result: the address is not taken | `s8c-recvescapesnamedresult` | build fails (CS1929) | `s8ptrmethodnamedresult` (receiver not used as a value) | `e2784f7e02` | mathutil, 4 sites |
| S9 | A type alias whose name is a C# keyword (`long`, `ulong`, `byte`) | `s9aliaskeywordname` | build fails (parse errors) | `s9ctl-aliasplainname` | `8638868a41` | libc (`long`, `ulong`), lib (`byte`) |
| S10 | A std identifier go2cs renames (`time.Location`, `Local`, `UTC`, `Second`), reached through a renamed import (`gotime "time"`) | `s10aliasedimportrenamed` | build fails (CS0426) | `s10ctl-plainimport` | `4b33aa9fad` | libc, 9 sites |
| S12 | Another package's alias named by qualifier (`stdio.FILE`, `type FILE = X_iobuf`) | `s12foreignaliasqualified` | build fails (CS0426) | `s12ctl-foreignnamed` | `2392980256` | libc, 1 site |
| S13 | A function ending in a terminating switch whose default falls through: C# sees a reachable end | `s13terminatingfallthrough` | build fails (CS0161) | `s13ctl-terminatingswitch` | `0ed7e0fc09` | libc pthread, 3 functions |
| S14 | A parallel assignment through dereferenced pointer casts (a swap) | `s14derefswap` | build fails (CS0131) | `s14ctl-derefswaptemp` | `0ed7e0fc09` | libc qsort, 3 swaps |
| S15 | An alias of an alias of an anonymous struct in one package, the alias declared before and in another file from its target | `s15b-aliasforwardotherfile` | build fails (a `/` in the alias target) | `s15ctl-aliasanonstruct` | `cf8f2cc8bb` | lib, 3,330 sites |
| S16 | An exported field whose type is an unexported alias of an anonymous struct, reached only as that field. go2cs's source generator chooses each declaration's accessibility from its Go name, so changing the C# by hand cannot fix it | `s16b-fieldonlyuse` | build fails (CS0050) | `s16ctl-exportedfieldexportedanon` | `518eae869d` | lib, 8 types and the types they reach |
| S17 | A field store through an explicit dereference of a pointer conversion, `(*(*T)(p)).F = x` | `s17explicitderefcastfieldstore` | build fails (CS0131) | `s17ctl-implicitderefcastfieldstore` (`(*T)(p).F = x`) | `b72be15017` | lib, 929 sites |
| S18 | Slicing a pointer cast whose C# type name holds a non-ASCII character: one character of the cast is lost | `s18nonasciiptrcastslice` | build fails (CS1003) | `s5ptrnamedarrayslice` (an ASCII name) | `e797b6716d` | the driver, 12 sites (`Δlibc.RawMem`, from go2cs's renamed import) |
| S19 | An imported package whose name differs from its directory (`modernc.org/sqlite/lib` is `package sqlite3`): references are spelled with the directory name | `s19pkgnamenotdirname` | build fails (CS0246) | `s19ctl-pkgnameisdirname` | `2392980256` | the driver, 2 sites |
| S20 | Slicing a pointer-to-array cast written in extra parentheses, `((*T)(p))[:n]`, or through an explicit dereference, `(*(*T)(p))[:n]`: the whole array is read as a value | `s20parenptrcastslice`, `s20b-derefptrcastslice` | runs, then crashes (AccessViolationException) | `s20ctl-ptrcastslice` | `b72be15017` | memory (Calloc), libc (5 sites); with a very large array type it overflows at start-up |
| S21 | A function implemented in assembly (a `.s` file) in a module outside the Go standard library: go2cs converts no assembly, so the call throws | `s21asmfuncinmodule` | runs, then `NotImplementedException` at the call | `s21ctl-gofuncinmodule` (the same function in Go) | `00176aac4a` | golang.org/x/sys/cpu `cpuid` and `xgetbv`, called by its `init` in every program importing modernc.org/libc |
| S22 | One struct read as another through `unsafe.Pointer` where both hold fixed-size arrays (ccgo's alignment wrapper). The go2cs runtime refuses it: a C# array field is a reference, where Go stores the elements inline | `s22reinterpretarraystruct` | runs, then panics: `cannot read native memory as struct {...}: its CLR layout is not its Go layout` | `s22ctl-reinterpretscalarstruct` (scalar fields) | `fcad1dc9a8` (runtime), `d770071de7` (converter) | lib, at start-up (`_sqlite3DigitPairs`); a runtime design limit, like S1 and S3 |
| S23 | An alias of the same anonymous struct, declared separately in two packages, is lifted to a separate C# type in each; Go sees one type, so a function stored as a number in one package and recovered with the other's signature no longer matches | go2cs test `ImportedAnonStructLiftReuse` | | | `d21fa11545` | lib's `_winHandleSeek` (from `_winSeekFile`) calls libc's `XSetFilePointerEx` through a function pointer typed with lib's `TLARGE_INTEGER`; libc declares the parameter as its own `largeInteger` |

Rows S2, S9, S12, S15, S16 and S19 are one family (alias and import naming), fixed together on go2cs master; S5 shares
S3's sites; S17 and S20 share an ingredient (an explicit dereference of a pointer conversion) and one fix. Variants that
did NOT reproduce their class are kept as negative evidence: `s8b-twotypessamemethod`, `s15aliaschainanonstruct`,
`s16exportedfieldunexportedanon`, `s17b-explicitderefplainptr`, `s18b-foreignnamedarrayslice`,
`s19b-namedstructrenamedimport`, `s19c-childpkgrenamedimport`, `s19d-aliasanoncastlocal`.

## Walls met running the demo

With rows S1-S23 fixed, the converted demo was built and run, and each further failure was recorded as a wall. A wall
is not reduced to a repro here; go2cs's own tests carry the reduction with each fix. Status is at go2cs master
`0c1d295ea8`. "Stand-in" means the demo runs looked past the wall with a hand rewrite of the converted C#, which is not
published.

| Wall | What it is | Status |
|---|---|---|
| W1 | `golang.org/x/sys/windows` refused: a pair of structs, one held inline and one not | fixed on master |
| W2 | SQLite's `lib` aliases emitted as raw Go: a 64 KB line reader, and a forward alias | fixed on master |
| W3 (S8) | a direct pointer method on a named result (`mathutil`) | fixed on master |
| W4 (S10) | libc's `time` members keep their Go names through a renamed import | fixed on master |
| W5a | an element address names another package's anonymous struct type | fixed on master |
| W5b | libc's `(*RawMem)(p)[:]` | fixed on master |
| W6 | an inline-layout attribute emitted without its `using` | fixed on master |
| W7 | an alias record keeps an unresolved marker (a forward alias chain) | fix written, not yet on master |
| W8 | SQLite's C string table read through a Go string header needs a base that never moves | stand-in; fix in progress |
| W9 | a C struct in native memory holding a Go `sync.Mutex` (SQLite's mutex on Windows, libc's pthread mutex on Linux) | stand-in; fix in progress |
| W9b | a write through `&(*T)(p).F` lands in a copy of the field | stand-in; fix in progress |
| W10 (S3) | `memcpy` copies into a temporary view | fixed on master |
| W11 | a renamed standard-library type spelled by its Go name in an interface record | stand-in; fix in progress |
| W12, W12b | an index, a struct-element read or a whole-array store through an array over C memory | stand-in at the latest run; fixes reached master afterwards, not yet re-measured |
| W13 | a renamed package variable's address spelled by its raw name | fix written, not yet on master |
| W14 | a large libc struct kept out of inline arrays reads null fields | fix written, not yet on master; the latest run needed no stand-in |
| W15 | a zero-length alignment field `[0][N]T` drops N (SQLite's bit vector, on every CREATE TABLE) | fixed on master (batch 309); the latest run needed no stand-in |
| L1 | an anonymous struct declared inside a function left unnamed | fixed on master |
| W17 | the array-over-C-memory fixes (W12's) name two anonymous structs declared inside a function before they are lifted | at master the conversion refuses here (exit code 1); fix in progress |
