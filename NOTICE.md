# Notices

This repository will carry C# converted by go2cs from the Go modules below. The converted code is a derivative of
those modules, and each module's own notices apply to it. They are kept verbatim in [LICENSES/](LICENSES/), named
`<module>--<file>`.

| Code | Notices |
|---|---|
| modernc.org/sqlite v1.62.0 (the driver and `lib`) | `modernc.org-sqlite--LICENSE` (BSD-3-Clause), `--AUTHORS`, `--LICENSE-3RD-PARTY.md` |
| SQLite 3.54.1, as transpiled into modernc.org/sqlite/lib | `modernc.org-sqlite--LICENSE-SQLITE` (public-domain dedication) |
| sqlite-vec, as bundled in modernc.org/sqlite | `modernc.org-sqlite--LICENSE-SQLITE_VEC` (MIT) |
| modernc.org/libc v1.80.1 | `modernc.org-libc--LICENSE`, `--AUTHORS`, `--LICENSE-3RD-PARTY.md` |
| modernc.org/memory v1.12.1 | `modernc.org-memory--LICENSE`, `--LICENSE-GO`, `--LICENSE-MMAP-GO`, `--LICENSE-LOGO`, `--AUTHORS` |
| modernc.org/mathutil v1.7.1 | `modernc.org-mathutil--LICENSE`, `--AUTHORS` |
| modernc.org/fileutil v1.4.0 | `modernc.org-fileutil--LICENSE`, `--AUTHORS` |
| golang.org/x/sys v0.49.0 | `golang.org-x-sys--LICENSE`, `--PATENTS` |
| github.com/dustin/go-humanize v1.1.0 | `github.com-dustin-go-humanize--LICENSE` |
| github.com/mattn/go-isatty v0.0.24 | `github.com-mattn-go-isatty--LICENSE` |
| github.com/ncruces/go-strftime v1.1.0 | `github.com-ncruces-go-strftime--LICENSE` |
| github.com/remyoudompheng/bigfft (2023-01-29) | `github.com-remyoudompheng-bigfft--LICENSE` |
| github.com/google/uuid v1.6.0 (linux targets) | `github.com-google-uuid--LICENSE` |

The files are copied byte for byte from each module as `go mod download` fetched it, through proxy.golang.org with
checksums verified by sum.golang.org.

This repository's own files (the demo, scripts and documentation) are under the MIT License in [LICENSE](LICENSE).
