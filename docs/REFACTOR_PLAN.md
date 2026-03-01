# PLAN MASTER: `driver/` Package — `github.com/cdvelop/sqlite-wasm`

## One-time Setup (run once before any phase)

```bash
# Install gotest — required for all test runs in this project
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

> Always use `gotest` instead of `go test`. It handles `-vet`, `-race`, `-cover`,
> WASM tests, and README badges automatically.

> **Repo:** `github.com/cdvelop/sqlite-wasm` (development sandbox)
> **Merge target:** `github.com/tinywasm/sqlite/driver/`

---

## Purpose

This repo is the **isolated workspace** to build a self-contained `driver/` package
that encapsulates the full SQLite engine (currently `modernc.org/sqlite` source).

**End goal:** copy the `driver/` folder verbatim into `tinywasm/sqlite/driver/`
with a single `sed` import-path substitution. After the copy, the adapter imports
only `_ "github.com/tinywasm/sqlite/driver"` and has zero `modernc.org/*` dependencies.

---

## Migration Strategy

During development, internal `driver/` imports reference this sandbox:

```
github.com/cdvelop/sqlite-wasm/driver/lib
github.com/cdvelop/sqlite-wasm/driver/vfs
github.com/cdvelop/sqlite-wasm/driver/vtab
```

At migration time, one `sed` switches to the final location:

```bash
sed -i 's|github.com/cdvelop/sqlite-wasm/driver|github.com/tinywasm/sqlite/driver|g' \
    driver/*.go driver/lib/*.go driver/vfs/*.go driver/vtab/*.go
```

The full migration procedure is in [8_MIGRATION.md](8_MIGRATION.md).

---

## Phases

| File | Phase | Goal | Status |
|------|-------|------|--------|
| [0_MODULE_CLEANUP.md](0_MODULE_CLEANUP.md) | 0 | Fix `go.mod`; remove package conflict; `go build ./...` passes | 🔲 Pending |
| [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md) | 1 | Move all engine source into `driver/`; rename packages; update imports | 🔲 Pending |
| [2_TESTS_PASS.md](2_TESTS_PASS.md) | 2 | Existing tests pass with new `driver/` structure (`gotest`) | 🔲 Pending |
| [3_TESTS_MOVE.md](3_TESTS_MOVE.md) | 3 | Move all tests to `tests/`; add build tags; `gotest` passes | 🔲 Pending |
| [4_TESTS_DOMAIN.md](4_TESTS_DOMAIN.md) | 4 | Subdivide tests by domain (conn, stmt, vfs, vtab, backup); coverage ≥ 90% | 🔲 Pending |
| [5_DEPS_SMALL.md](5_DEPS_SMALL.md) | 5 | Inline `modernc.org/mathutil`, `modernc.org/fileutil`, `modernc.org/memory` into `driver/` | 🔲 Pending |
| [6_DEPS_LIBC.md](6_DEPS_LIBC.md) | 6 | Inline `modernc.org/libc` into `driver/` — zero `modernc.org/*` in `go.mod` | 🔲 Pending |
| [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md) | 7 | Replace `google/uuid` → `tinywasm/unixid`; clean `go.mod` to only `tinywasm/*` and stdlib deps | 🔲 Pending |
| [8_MIGRATION.md](8_MIGRATION.md) | 8 | Write & validate `scripts/migrate_to_tinywasm.sh`; dry-run copy into `tinywasm/sqlite` | 🔲 Pending |

---

## Target Architecture

```
cdvelop/sqlite-wasm/                   ← this repo (sandbox)
├── go.mod                             ← module github.com/cdvelop/sqlite-wasm
├── driver/                            ← THE deliverable (copied to tinywasm/sqlite/driver/)
│   ├── sqlite.go                      ← engine init() + driver registration (package driver)
│   ├── conn.go                        ← connection handling
│   ├── stmt.go                        ← statement handling
│   ├── rows.go                        ← row iteration
│   ├── tx.go                          ← transactions
│   ├── result.go                      ← sql.Result implementation
│   ├── error.go                       ← error types
│   ├── convert.go                     ← type conversion
│   ├── vtab.go                        ← virtual tables (flat in driver/)
│   ├── backup.go                      ← backup API
│   ├── ...                            ← other engine files (all flat in driver/)
│   ├── lib/                           ← package sqlite3 (C→Go transpiled, auto-generated)
│   │   ├── sqlite_linux_amd64.go
│   │   └── ...
│   ├── vfs/                           ← package vfs (24 files, auto-generated VFS layer)
│   │   └── ...
│   └── vtab/                          ← package vtab (2 files, virtual table helpers)
│       └── ...
├── tests/                             ← all tests (Phase 3+)
│   ├── conn_test.go
│   ├── stmt_test.go
│   ├── vfs_test.go
│   └── ...
└── scripts/
    └── migrate_to_tinywasm.sh         ← Phase 8: copies driver/ + fixes import paths
```

### Sub-package Policy

| Sub-package | Reason for sub-package |
|-------------|------------------------|
| `driver/lib/` | `package sqlite3` — auto-generated C→Go; cannot merge |
| `driver/vfs/` | `package vfs` — 24 auto-generated platform files; cannot merge |
| `driver/vtab/` | `package vtab` — internal vtab helpers imported by `driver/vtab.go` |

---

## Development Rules

- **`gotest`** for all test runs (not `go test`). See One-time Setup above.
- **Max 500 lines per file** (applies to new hand-written files, not auto-generated sources).
- **No external assertion libraries.** Standard `testing` only.
- **No global state.** Dependency injection via interfaces.
- Coverage target: **≥ 90%** from Phase 4 onward.

---

## References

- [`cdvelop/sqlite-wasm`](https://github.com/cdvelop/sqlite-wasm) — this repo (sandbox)
- [`tinywasm/sqlite`](https://github.com/tinywasm/sqlite) — merge target (adapter)
- [`tinywasm/orm`](https://github.com/tinywasm/orm) — ORM abstraction layer
- [`tinywasm/unixid`](https://github.com/tinywasm/unixid) — replaces `google/uuid`
- [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) — original engine source
