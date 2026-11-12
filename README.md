# sqlite-cs

[![converted by go2cs](https://go2cs.net/images/go2cs-badge.svg)](https://go2cs.net)

[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite), the pure-Go SQLite, converted to C# by
[go2cs](https://github.com/ritchiecarroll/go2cs), and run against its own tests.

> [!IMPORTANT]
> **This is a work in progress, and a proof of what go2cs can do. It is not a SQLite library for .NET.** For that, use
> [Microsoft.Data.Sqlite](https://learn.microsoft.com/dotnet/standard/data/sqlite/) or another maintained package.
> This repository is not affiliated with or endorsed by the SQLite project or by the authors of modernc.org/sqlite.

## Status

The status below is measured, and it changes as the work moves. It is stated as it is today, failures included.

| Stage | windows/amd64 |
|---|---|
| Native Go build and its own tests | **pass**: 432 tests pass, 4 skip, 0 fail (`go test -count=1 .`, 701 s) |
| Conversion to C# | **refused**: go2cs stops because some types come out as raw Go text (rows S1, S2) |
| C# build | **fails**. With each failure class corrected by hand for diagnosis, the whole demo builds |
| The demo, converted | built only with those hand corrections; it **stops at start-up** (row S22) |
| Its own tests, converted | not reached |

The first target is windows/amd64: its dependency closure has no cgo and a single assembly file. Other targets follow
once the first one runs.

### Where the build stops

To see past each failure, the converted C# was corrected by hand, one class at a time, and rebuilt. Those corrections are
for diagnosis only and are not in this repository. Each class found is reduced to a one-file reproduction in
[repros/](repros/), with a table of what each does today.

- `github.com/remyoudompheng/bigfft`, `modernc.org/memory`, `github.com/ncruces/go-strftime`, `modernc.org/mathutil`
  and `modernc.org/libc` build once rows S3-S14 are corrected.
- `modernc.org/sqlite/lib` (SQLite itself) then builds once rows S15-S17 are corrected, and the driver
  (`modernc.org/sqlite`) once rows S12, S18 and S19 are.
- The demo then builds. Run, it stops during start-up, before any SQL: first at row S20 (memory's allocator and libc),
  then at row S21 (`golang.org/x/sys/cpu`), and now at row S22. SQLite's own initialization reads one struct holding a
  fixed-size array as another, and the go2cs runtime refuses that, loudly, because a C# array field is a reference
  where Go stores the elements inline.
- Behind that, row S1 is still ahead: SQLite calls C function pointers by turning numbers back into Go functions,
  which go2cs cannot yet do.

## What is converted

| Module | Version | Role |
|---|---|---|
| modernc.org/sqlite | v1.62.0 | the driver, and `lib`: SQLite 3.54.1 transpiled from C to Go by ccgo |
| modernc.org/libc | v1.80.1 | C's runtime for the transpiled code: memory, strings, I/O |
| modernc.org/memory | v1.12.1 | libc's allocator |
| modernc.org/mathutil | v1.7.1 | math helpers |
| modernc.org/fileutil | v1.4.0 | file helpers |
| golang.org/x/sys | v0.49.0 | system calls |
| github.com/dustin/go-humanize | v1.1.0 | formatting (tests) |
| github.com/mattn/go-isatty | v0.0.24 | terminal detection |
| github.com/ncruces/go-strftime | v1.1.0 | time formatting |
| github.com/remyoudompheng/bigfft | 2023-01-29 | big-number multiplication |

On windows/amd64 the closure is 23 packages and about 407,000 lines of Go (14 MB); `modernc.org/sqlite/lib` alone is
331,000 lines. The Go standard library is not converted here: go2cs ships it already converted.

## Why it is hard

The transpiled SQLite keeps C's memory model inside Go: pointers are `uintptr` numbers into libc's own heap, every field
access is a `*(*T)(unsafe.Pointer(p + offset))`, and function pointers are stored as numbers and turned back into
callable functions. The windows/amd64 closure has about 40,000 `unsafe.Pointer` sites. .NET's garbage collector moves
objects and has no "address of a function", so this is where a Go-to-C# conversion is tested hardest. Each failure
class found here is recorded and fixed in go2cs where it can be.

## Reproduce

Requirements: Go 1.26.10, the .NET 10 SDK, and a clone of go2cs.

```
# the native oracle
cd demo
go run .
# expected:
#   sqlite 3.54.1
#   2 beta 2.25
#   1 alpha 1.5
#   3 gamma -3
```

The demo opens an in-memory database through `database/sql`, creates a table, inserts three rows and queries them. It
is not part of modernc.org/sqlite. The conversion commands are added here when the conversion runs.

## Licenses

- The converted code derives from modernc.org/sqlite (BSD-3-Clause), SQLite itself (public domain), sqlite-vec (MIT),
  modernc.org/libc and its third-party sources, and the other modules above. Their notices are kept verbatim in
  [LICENSES/](LICENSES/); [NOTICE.md](NOTICE.md) says which notice covers which code.
- This repository's own files (the demo, scripts and documentation) are under the MIT License in [LICENSE](LICENSE).
- go2cs, the converter, is licensed separately; see its repository.
