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
| S1 | A C function pointer kept as a `uintptr` and called by rebuilding a Go func value from `struct{ uintptr }`; the converter spells the reinterpret's source type as raw Go | `s1funcptrstruct` | conversion refused: `struct{uintptr}` emitted as raw Go | | open: the design is approved | 503 sites (lib, libc). Behind it, a delegate cannot be rebuilt from a number: a runtime design gap |
| S2 | An alias in another package of an alias of an anonymous struct | `s2aliasanonstruct` | conversion refused: the struct emitted as raw Go | `s2ctl-namedstruct` | `cf8f2cc8bb` | 4 sites (driver's vtab.go, pagecache.go) |
| S3 | A slice of a pointer-to-array cast COPIES the memory it views, so a write through it is lost | `s3ptrarrayviewcopy` | runs, wrong output: `8 [0 0 0 0 0 0 0 0 0] 4 0` (Go `8 [1 2 3 4 5 6 7 8 0] 4 4`) | | open: with the pointer-view work | 20 sites, the converter warns at each; the allocator's realloc (memory.go:502) among them |
| S4 | A defined type over another defined slice type (`type fermat nat`, `type nat []big.Word`) loses slice operations and its conversion to `nat` | `s4namedslicechain` | build fails (CS1503, CS0021, ...) | `s4ctl-namedslice` | open: routed | bigfft, 131 errors |
| S5 | Slicing a pointer to a NAMED array type (`type rawmem [N]byte`) gives a slice of the array type, not of `byte` | `s5ptrnamedarrayslice` | runs, then `ArgumentException: Cannot use type 'rawmem'` | | `b72be15017` (the element type) | memory (5 sites), libc's RawMem (about 9 sites); the same sites as S3 |
| S6 | `literalErr(b)`, a byte converted to a named string type (Go's `string(rune(b))`) | `s6bytetonamedstring` | build fails (CS0030) | `s6ctl-stringtonamedstring` | open: routed | go-strftime, 1 site |
| S7 | A user struct type named `float`: its composite literal is emitted as `float64` | `s7typenamedfloat` | build fails (CS1739, CS1929) | `s7ctl-typenamedfloatx` | open: routed | mathutil, 1 site |
| S8 | A pointer method that uses its receiver as a value (`fmt.Errorf("%T", x)`), called on an addressable named result: the address is not taken | `s8c-recvescapesnamedresult` | build fails (CS1929) | `s8ptrmethodnamedresult` (receiver not used as a value) | open: routed | mathutil, 4 sites |
| S9 | A type alias whose name is a C# keyword (`long`, `ulong`, `byte`) | `s9aliaskeywordname` | build fails (parse errors) | `s9ctl-aliasplainname` | `8638868a41` | libc (`long`, `ulong`), lib (`byte`) |
| S10 | A std identifier go2cs renames (`time.Location`, `Local`, `UTC`, `Second`), reached through a renamed import (`gotime "time"`) | `s10aliasedimportrenamed` | build fails (CS0426) | `s10ctl-plainimport` | open: a fix is decided, not yet on master | libc, 9 sites |
| S12 | Another package's alias named by qualifier (`stdio.FILE`, `type FILE = X_iobuf`) | `s12foreignaliasqualified` | build fails (CS0426) | `s12ctl-foreignnamed` | `2392980256` | libc, 1 site |
| S13 | A function ending in a terminating switch whose default falls through: C# sees a reachable end | `s13terminatingfallthrough` | build fails (CS0161) | `s13ctl-terminatingswitch` | open: routed | libc pthread, 3 functions |
| S14 | A parallel assignment through dereferenced pointer casts (a swap) | `s14derefswap` | build fails (CS0131) | `s14ctl-derefswaptemp` | open: routed | libc qsort, 3 swaps |
| S15 | An alias of an alias of an anonymous struct in one package, the alias declared before and in another file from its target | `s15b-aliasforwardotherfile` | build fails (a `/` in the alias target) | `s15ctl-aliasanonstruct` | `cf8f2cc8bb` | lib, 3,330 sites |
| S16 | An exported field whose type is an unexported alias of an anonymous struct, reached only as that field. go2cs's source generator chooses each declaration's accessibility from its Go name, so changing the C# by hand cannot fix it | `s16b-fieldonlyuse` | build fails (CS0050) | `s16ctl-exportedfieldexportedanon` | `518eae869d` | lib, 8 types and the types they reach |
| S17 | A field store through an explicit dereference of a pointer conversion, `(*(*T)(p)).F = x` | `s17explicitderefcastfieldstore` | build fails (CS0131) | `s17ctl-implicitderefcastfieldstore` (`(*T)(p).F = x`) | `b72be15017` | lib, 929 sites |
| S18 | Slicing a pointer cast whose C# type name holds a non-ASCII character: one character of the cast is lost | `s18nonasciiptrcastslice` | build fails (CS1003) | `s5ptrnamedarrayslice` (an ASCII name) | `e797b6716d` | the driver, 12 sites (`Δlibc.RawMem`, from go2cs's renamed import) |
| S19 | An imported package whose name differs from its directory (`modernc.org/sqlite/lib` is `package sqlite3`): references are spelled with the directory name | `s19pkgnamenotdirname` | build fails (CS0246) | `s19ctl-pkgnameisdirname` | `2392980256` | the driver, 2 sites |
| S20 | Slicing a pointer-to-array cast written in extra parentheses, `((*T)(p))[:n]`, or through an explicit dereference, `(*(*T)(p))[:n]`: the whole array is read as a value | `s20parenptrcastslice`, `s20b-derefptrcastslice` | runs, then crashes (AccessViolationException) | `s20ctl-ptrcastslice` | `b72be15017` | memory (Calloc), libc (5 sites); with a very large array type it overflows at start-up |
| S21 | A function implemented in assembly (a `.s` file) in a module outside the Go standard library: go2cs converts no assembly, so the call throws | `s21asmfuncinmodule` | runs, then `NotImplementedException` at the call | `s21ctl-gofuncinmodule` (the same function in Go) | open: designed with S1 | golang.org/x/sys/cpu `cpuid` and `xgetbv`, called by its `init` in every program importing modernc.org/libc |
| S22 | One struct read as another through `unsafe.Pointer` where both hold fixed-size arrays (ccgo's alignment wrapper). The go2cs runtime refuses it: a C# array field is a reference, where Go stores the elements inline | `s22reinterpretarraystruct` | runs, then panics: `cannot read native memory as struct {...}: its CLR layout is not its Go layout` | `s22ctl-reinterpretscalarstruct` (scalar fields) | open: in progress (inline fixed-size arrays) | lib, at start-up (`_sqlite3DigitPairs`); a runtime design limit, like S1 and S3 |
| S23 | An alias of the same anonymous struct, declared separately in two packages, is lifted to a separate C# type in each; Go sees one type, so a function stored as a number in one package and recovered with the other's signature no longer matches | not yet reduced | | | open: routed | lib's `_winHandleSeek` (from `_winSeekFile`) calls libc's `XSetFilePointerEx` through a function pointer typed with lib's `TLARGE_INTEGER`; libc declares the parameter as its own `largeInteger` |

Rows S2, S9, S12, S15, S16 and S19 are one family (alias and import naming), fixed together on go2cs master; S5 shares
S3's sites; S17 and S20 share an ingredient (an explicit dereference of a pointer conversion) and one fix. Variants that
did NOT reproduce their class are kept as negative evidence: `s8b-twotypessamemethod`, `s15aliaschainanonstruct`,
`s16exportedfieldunexportedanon`, `s17b-explicitderefplainptr`, `s18b-foreignnamedarrayslice`,
`s19b-namedstructrenamedimport`, `s19c-childpkgrenamedimport`, `s19d-aliasanoncastlocal`.
