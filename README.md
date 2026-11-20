# sqlite-cs

[![converted by go2cs](https://cdn.jsdelivr.net/gh/ritchiecarroll/go2cs@master/docs/images/go2cs-badge.svg)](https://go2cs.net)

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
| Native Go build and its own tests | **pass**: 432 tests pass, 4 skip, 0 fail (`go test -count=1 .`, Go 1.26.10, 701 s) | **pass**: 438 tests pass, 3 skip, 0 fail (`go test -count=1 .`, Go 1.26.11, 495 s) |
| Conversion to C# | **completes** (`go2cs -recurse=nuget`) | **completes** (`go2cs -recurse=nuget`) |
| C# build | **builds**, no stand-ins (go2cs 1.26.11.1, the published standard library) | **builds**, no stand-ins (go2cs 1.26.11.1, the published standard library) |
| The demo, converted | output **byte-identical** to native Go, no stand-ins | output **byte-identical** to native Go, no stand-ins |
| The file-backed demo, converted | output **byte-identical** to native Go in both journal modes, no stand-ins | output **byte-identical** to native Go in both journal modes, no stand-ins |
| SQLite's own tests, converted | **not yet run** | **in progress**: they need go2cs fixes newer than 1.26.11.1 |

The demo runs were measured on 2026-11-20 with Go 1.26.11, the .NET 10 SDK and the go2cs 1.26.11.1 release,
unmodified, comparing the converted program's output with the native program's using `cmp`. The Linux runs are on
Ubuntu 24.04 (x64). On both operating systems the standard library comes from the `go.*` 1.26.11.1 packages on
nuget.org, with no local package feed and no go2cs clone.

In the file-backed demo's WAL pass on Linux (traced on 2026-11-15), the converted program opens, maps and locks
SQLite's shared-memory file (`-shm`) exactly as the native program does, call for call, as traced with `strace`:
SQLite's WAL index lives in a shared memory mapping in the converted program too.

A stand-in is a hand rewrite of the converted C#, made only to see past a wall that go2cs did not yet handle. Earlier
runs needed stand-ins; the conversion published here needs none, on Windows or on Linux.
[repros/](repros/README.md#walls-met-running-the-demo) lists every wall and where it stands.

The converted C# is published in this repository in [csharp/](csharp/), one tree per target and demo;
[csharp/PROVENANCE.md](csharp/PROVENANCE.md) records the converter, the command and the environment that produced it. Every file's SHA-256 is in
[csharp/MANIFEST.sha256](csharp/MANIFEST.sha256); check it from the repository root with
`sha256sum -c csharp/MANIFEST.sha256`.

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

With those fixed, the demo was run, and each further failure became a wall, W1-W18 and L1. The go2cs 1.26.11.1 release
fixes them, except W18 on Windows, which the `-inline-arrays` option under [Reproduce](#reproduce) answers.
[repros/](repros/README.md#walls-met-running-the-demo) lists them.

<!-- sqlite-cs:performance -->
## Performance

The converted SQLite is much slower than native Go today. That is expected for code that was C, then Go, and is now
C# on a different runtime, and it is improving. The table gives seconds. Each time is split into start-up (the time
before any SQLite work begins) and work (everything after it), so the work column can be compared directly.

<!-- Measured by G (go2cs fleet) 2026-11-20, COORD item 4. The workload (a fixed program, not part of
     modernc.org/sqlite): one database FILE in a fresh temporary directory; open it and create
     t(id INTEGER PRIMARY KEY, k INTEGER, v TEXT, f REAL); insert 100,000 rows in ONE transaction through one prepared
     statement; CREATE INDEX on k; one GROUP BY aggregate over every row (100 groups), then 10,000 point lookups through
     the index, each a COUNT and SUM over 100 rows. SQLite's default journal mode. Stdout is three result lines, equal
     by SHA-256 between native and converted on every case before any run was timed; each phase's time goes to stderr.
     Start-up = the median wall of the same program on an empty table (the same code with 0 rows); work = the median
     wall minus that start-up; median of 5 timed runs after 1 warm-up, one process per run. Phase times are the
     program's own, medians of the same 5 runs.
     Converted with go2cs 1.26.11.1 (commit 5dd5cafd5b): go2cs -recurse=nuget
     -inline-arrays=modernc.org/libc#largeInteger <program> <out>, with GO2CSPATH set to an empty folder and no
     go2cs folder in the user profile; built Release against go.* 1.26.11.1 from nuget.org, .NET SDK 10.0.400. The
     closure (23 packages, modernc.org/sqlite v1.62.0, modernc.org/libc v1.80.1) built with 0 errors and NO hand
     stand-ins. Traces: two 20,000-row runs under dotnet-trace (sampled thread time + GC events; GC verbose for
     allocation by type: 31% and 34% of the span in GC pauses, 10.6 GB allocated in each). -->
Measured on G-LAPTOP (AMD Ryzen 5 PRO 6650U, 30.8 GB of RAM, Windows 11), 2026-11-20. The native column is Go 1.26.11
for windows/amd64 with modernc.org/sqlite v1.62.0; the converted column is win-x64 on .NET 10 (SDK 10.0.400, runtime
10.0.11, the default workstation garbage collector) against go.* 1.26.11.1.

The workload is a fixed program over one database file: it inserts 100,000 rows in one transaction, indexes them,
and queries them back (one aggregate over every row, then 10,000 lookups through the index).

| Step | Go | Converted, JIT |
|:--|--:|--:|
| Start-up (the same program on an empty table) | 0.04 | 2.41 |
| Open the file and create the table | 0.01 | 0.95 |
| Insert 100,000 rows in one transaction | 0.17 | 40.71 |
| Create the index | 0.04 | 33.46 |
| Query: the aggregate and 10,000 indexed lookups | 0.87 | 152.53 |
| **Total**, start-up + work | **0.04 + 1.08 = 1.12** | **2.41 + 227.47 = 229.88** |

Every converted run printed exactly what the native build printed: the output was compared by SHA-256 before any run
was timed. Peak memory was 23 MB for Go on windows/amd64 and 269 MB for the converted program on win-x64.

### Where the time goes

These are the largest differences we measured on win-x64, and each is a target for future work in go2cs.

- **Allocation and garbage collection.** A 20,000-row run allocated 10.6 GB, against a peak heap of 39 MB: short-lived
  objects that stand for C's raw memory. Unsafe pointer values were 39% of the bytes, arrays of `uintptr` 21%, and
  typed views of C structs in native memory about 30%. The run made 1,303 collections, and their pauses were about a
  third of its time.
- **Start-up.** The converted program spends 2.4 s before it does any work, where native Go on windows/amd64 needs
  0.04 s: compiling methods on first call and running the initializers of every package in the closure.

Where the rest of the time goes is not yet measured: the sampling profiler used here attributes time only at safe
points, and a first fix aimed at what it showed did not change the run time.
<!-- Corrected 2026-11-20: an earlier version of this list named a contended lock at the native boundary (80% of the
     working thread, split 42% / 24%) and a busy finalizer thread. Both readings came from a thread-time sampler that
     stops threads only at safe points; a fix that removed the first lock site moved the sampled wait elsewhere and left
     the run time unchanged (A 13.0 / 13.7 / 12.7 s, B 13.3 / 13.2 / 13.5 s, 20,000 rows), so they were withdrawn. The
     garbage-collection figures come from runtime GC events and stand. -->
<!-- /sqlite-cs:performance -->

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
331,000 lines. The Go standard library is not converted here: go2cs ships it already converted. On linux/amd64,
`github.com/google/uuid` is referenced as its published NuGet package rather than converted.

## What makes this conversion hard

The transpiled SQLite keeps C's memory model inside Go: pointers are `uintptr` numbers into libc's own heap, every field
access is a `*(*T)(unsafe.Pointer(p + offset))`, and function pointers are stored as numbers and turned back into
callable functions. The windows/amd64 closure has about 40,000 `unsafe.Pointer` sites. .NET's garbage collector moves
objects and has no "address of a function", so this is where a Go-to-C# conversion is tested hardest. Each failure
class found here is recorded and fixed in go2cs where it can be.

## Reproduce

Requirements: Go 1.26.11, the .NET 10 SDK, and the go2cs tool 1.26.11.1 (`dotnet tool install -g go2cs`). No clone of
go2cs is needed.

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

To convert it with its whole module closure, run go2cs from inside the demo's folder. `<demo>` is this repository's
`demo/` folder and `<out>` is an empty output folder. Always give `<out>`: without it, go2cs writes beside its input.

```
cd <demo>
go2cs -recurse=nuget -inline-arrays=modernc.org/libc#largeInteger <demo> <out>
```

On Windows, libc's `largeInteger` must be an inline array on both sides of the libc/sqlite3 boundary, which the
`-inline-arrays` option asks for; a later go2cs derives it by itself. On Linux the converted C# is the same with or
without it.

Build the converted demo. The standard library is restored from nuget.org. The generated `Directory.Build.props`
already selects the target the tree was converted for, so the same command builds on Windows and on Linux:

```
dotnet build <out>/src/sqlitecs/demo/sqlitecs.demo.csproj -c Release
```

Run the program the build writes under `<out>/.artifacts/bin/`, then compare its output with the native program's:

```
dotnet <out>/.artifacts/bin/<...>/Release/net10.0/<program>.dll > converted.txt
(cd <demo> && go run .) > native.txt
cmp converted.txt native.txt
```

`demo/file/` is a second demo, its own Go module, over a database file. In SQLite's default journal mode and then in
WAL mode, it creates a database file in a new temporary directory, adds a table, an index and three rows, and closes
it; it then reopens the file, checks its integrity, queries it and removes it. Natively it prints 16 lines. Convert,
build and compare it the same way, with `demo/file/` as `<demo>`; its project is under `<out>/src/sqlitecs/demo/file/`.

The trees in [csharp/](csharp/) are the output of this conversion, with the exact command and environment recorded in
[csharp/PROVENANCE.md](csharp/PROVENANCE.md), so a tree can also be built and run directly without converting it
yourself.

## Licenses

- The converted code derives from modernc.org/sqlite (BSD-3-Clause), SQLite itself (public domain), sqlite-vec (MIT),
  modernc.org/libc and its third-party sources, and the other modules above. Their notices are kept verbatim in
  [LICENSES/](LICENSES/); [NOTICE.md](NOTICE.md) says which notice covers which code.
- This repository's own files (the demo, scripts and documentation) are under the MIT License in [LICENSE](LICENSE).
- go2cs, the converter, is licensed separately; see its repository.
