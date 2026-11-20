<!-- sqlite-cs:performance -->
# Performance

The converted SQLite is much slower than native Go today. That is expected for code that was C, then Go, and is now
C# on a different runtime, and it is improving. The table gives seconds. Each time is split into start-up (the time
before any SQLite work begins) and work (everything after it), so the work column can be compared directly.

<!-- Measured 2026-11-20. The workload (a fixed program, not part of
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
Measured 2026-11-20 on a laptop with an AMD Ryzen 5 PRO 6650U, 30.8 GB of RAM and Windows 11. The native column is Go 1.26.11
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

## Where the time goes

These are the largest differences we measured on win-x64, and each is a target for future work in go2cs.

- **Converting between unsafe pointers and numbers.** A CPU profile of a 20,000-row run, taken with Windows' kernel
  sampler, shows the working thread spending about half its time in the .NET runtime and 12% in go2cs's runtime
  library, against 11% in the converted SQLite code itself. The runtime is entered mostly from converting
  `unsafe.Pointer` values to and from `uintptr` numbers in SQLite's bytecode engine (`VdbeExec`), and those
  conversions create the short-lived objects behind the allocation below.
- **Allocation and garbage collection.** A 20,000-row run allocated 10.6 GB, against a peak heap of 39 MB: short-lived
  objects that stand for C's raw memory. Unsafe pointer values were 39% of the bytes, arrays of `uintptr` 21%, and
  typed views of C structs in native memory about 30%. The run made 1,303 collections, and their pauses were about a
  third of its time.
- **Start-up.** The converted program spends 2.4 s before it does any work, where native Go on windows/amd64 needs
  0.04 s: compiling methods on first call and running the initializers of every package in the closure.

<!-- Corrected 2026-11-20: an earlier version of this list named a contended lock at the native boundary (80% of the
     working thread, split 42% / 24%) and a busy finalizer thread. Both readings came from a thread-time sampler that
     stops threads only at safe points; a fix that removed the first lock site moved the sampled wait elsewhere and left
     the run time unchanged (A 13.0 / 13.7 / 12.7 s, B 13.3 / 13.2 / 13.5 s, 20,000 rows), so they were withdrawn. The
     garbage-collection figures come from runtime GC events and stand. A kernel CPU profile (Windows Performance
     Recorder, 1 ms samples, 29,211 for the process, CLR method events for managed frames) then showed no lock in the
     exclusive ranking (System.Collections.Concurrent 0.2%): work thread 19,451 samples = coreclr 49%, golib 12%, the
     converted SQLite 11%, JIT 6%, kernel 6%; a background thread spent 7,060 samples recompiling hot methods (tiering)
     on another core; the module initializer 2,192 at start-up. -->
<!-- /sqlite-cs:performance -->
