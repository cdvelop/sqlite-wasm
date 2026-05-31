# driver — SQLite Engine Package

This package is a self-contained embedding of the SQLite engine
(originally `modernc.org/sqlite`) for use as a `database/sql` driver.

## Registration

The `init()` function in `sqlite.go` registers the `"sqlite"` driver name
automatically on blank import:

    import _ "github.com/tinywasm/sqlite/driver"

## Sub-packages

| Package | Description |
|---------|-------------|
| `driver/lib/` | `package sqlite3` — auto-generated C→Go SQLite amalgamation |
| `driver/vfs/` | `package vfs` — virtual file system layer |
| `driver/vtab/` | `package vtab` — virtual table helpers |

## Dependencies

`modernc.org/libc` is an accepted direct dependency (CGo-free C stdlib emulation
required by the SQLite engine). All other deps are indirect entries managed
automatically by `go mod tidy`.

## Origin

Built and refined in `github.com/cdvelop/sqlite-wasm`.
