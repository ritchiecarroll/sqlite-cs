# sqlite-cs

[![converted by go2cs](https://go2cs.net/images/go2cs-badge.svg)](https://go2cs.net)

[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite), the pure-Go SQLite, converted to C# by [go2cs](https://github.com/ritchiecarroll/go2cs), and run against its own tests.

> [!IMPORTANT]
> **This is a work in progress, and a proof of what go2cs can do, i.e., produce operational .NET code transpiled from
> Go source. The code in this repository is _not_ intended to replace any existing SQLite libraries for .NET.** To use
> SQLite in .NET, use the official [Microsoft.Data.Sqlite](https://learn.microsoft.com/dotnet/standard/data/sqlite/) or
> another maintained package. This repository is not affiliated with or endorsed by the SQLite project or by the
> authors of modernc.org/sqlite.

## Status

The status below is measured, and it changes as the work moves. It is stated as it is today, failures included.

| Stage | windows/amd64 | linux/amd64 |
|---|---|---|
| Native Go build and its own tests | **pass**: 432 tests pass, 4 skip, 0 fail (`go test -count=1 .`, Go 1.26.10, 701 s) | not measured here |
| Conversion to C# | **completes** (`go2cs -recurse`) | **completes** (`go2cs -recurse`) |
| C# build | **builds**, with hand stand-ins for the walls go2cs has not fixed yet | **builds**, with hand stand-ins, fewer |
| The demo, converted | output **byte-identical** to native Go, with stand-ins | output **byte-identical** to native Go, with stand-ins |
| SQLite's own tests, converted | **not yet run** | **not yet run** |

The demo runs were measured on 2026-11-15 with Go 1.26.11 and the .NET 10 SDK, comparing the converted program's output
with the native program's using `cmp`. The Linux runs are on WSL2 (Ubuntu 22.04). The Windows run used an older go2cs
than the latest Linux run, so it needed more stand-ins, including a conversion option to look past wall W1, which go2cs
now fixes.

A stand-in is a hand rewrite of the converted C#, made only to see past a wall that go2cs does not yet handle. Stand-ins
are not in this repository, and each wall's real fix goes into go2cs. In the latest Linux run, the stand-ins left are
walls W8, W9, W9b, W11, W12 and W12b. Their fixes are in progress in go2cs. [repros/](repros/README.md#walls-met-running-the-demo)
lists every wall and where it stands.

The converted C# is published in this repository with the go2cs 1.26.11.1 release. Until then, convert it yourself with
the commands under [Reproduce](#reproduce).

Performance parity is not a goal of go2cs, nor of this proof (see the 
[FAQ](https://go2cs.net/FAQ.html#is-converted-code-slower-than-go)). go2cs strives to bring Go code to .NET with visual
and behavioral parity first; performance optimizations are currently secondary. Expect the converted SQLite code to
run slower in .NET than the native Go build.

### How it got here

The first conversion, on windows/amd64, was refused. Correcting the converted C# by hand past each failure, through the
build and then start-up, found the classes recorded as rows S1-S23. Each is reduced to a small reproduction in
[repros/](repros/), and every one is now fixed on go2cs's master branch. The hardest were C function
pointers kept as numbers (S1), assembly bodies outside the standard library (S21), and structs holding fixed-size
arrays read through `unsafe.Pointer` (S22).

With those fixed, the demo was run, and each further failure became a wall, W1-W15 and L1. Most are fixed in go2cs;
the rest are looked past with stand-ins, as above. [repros/](repros/README.md#walls-met-running-the-demo) lists both.

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

## What makes this conversion hard

The transpiled SQLite keeps C's memory model inside Go: pointers are `uintptr` numbers into libc's own heap, every field
access is a `*(*T)(unsafe.Pointer(p + offset))`, and function pointers are stored as numbers and turned back into
callable functions. The windows/amd64 closure has about 40,000 `unsafe.Pointer` sites. .NET's garbage collector moves
objects and has no "address of a function", so this is where a Go-to-C# conversion is tested hardest. Each failure
class found here is recorded and fixed in go2cs where it can be.

## Reproduce

Requirements: Go 1.26.11, the .NET 10 SDK, and a clone of [go2cs](https://github.com/ritchiecarroll/go2cs).

The demo opens an in-memory database through `database/sql`, creates a table, inserts three rows and queries them. It
is not part of modernc.org/sqlite.

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

To convert it, build go2cs from its clone, then convert the demo with its whole module closure. `<go2cs>` is the go2cs
clone, `<bin>` is where the go2cs binary goes, `<demo>` is this repository's `demo/` folder, and `<out>` is an empty
output folder. Always give `<out>`: without
it, go2cs writes beside its input.

```
go -C <go2cs>/src/go2cs build -o <bin>/go2cs .
cd <demo>
<bin>/go2cs -recurse -go2cspath <go2cs>/src <demo> <out>
```

Build the converted demo. On Linux, `-p:GoTargetOS=linux` is required; without it, the Windows flavour of the
converted standard library is linked.

```
# Windows
dotnet build <out>/src/sqlitecs/demo/sqlitecs.demo.csproj -c Release
# Linux
dotnet build <out>/src/sqlitecs/demo/sqlitecs.demo.csproj -c Release -p:GoTargetOS=linux
```

Run the program the build writes under `<out>/.artifacts/bin/`, then compare its output with the native program's:

```
dotnet <out>/.artifacts/bin/<...>/Release/net10.0/<program>.dll > converted.txt
(cd <demo> && go run .) > native.txt
cmp converted.txt native.txt
```

Until go2cs fixes the remaining walls, the converted C# needs the hand stand-ins described above before it builds and
runs. They are not published, so at go2cs master these steps stop at the first open wall.

## Licenses

- The converted code derives from modernc.org/sqlite (BSD-3-Clause), SQLite itself (public domain), sqlite-vec (MIT),
  modernc.org/libc and its third-party sources, and the other modules above. Their notices are kept verbatim in
  [LICENSES/](LICENSES/); [NOTICE.md](NOTICE.md) says which notice covers which code.
- This repository's own files (the demo, scripts and documentation) are under the MIT License in [LICENSE](LICENSE).
- go2cs, the converter, is licensed separately; see its repository.
