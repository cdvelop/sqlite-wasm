# Phase 1: Organize All Engine Source into `driver/`

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [0_MODULE_CLEANUP.md](0_MODULE_CLEANUP.md) ← must be complete
> **Next:** [2_TESTS_PASS.md](2_TESTS_PASS.md)

## Prerequisites

```bash
# Install gotest (required — do NOT use go test directly)
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Move all SQLite engine source files from the repo root into a `driver/` subdirectory.
Rename `package sqlite` → `package driver`. Update all internal import paths.

After this phase, the `driver/` folder is the **deliverable** — it can be copied
into `tinywasm/sqlite/driver/` with minimal changes.

---

## Sub-package Layout After This Phase

```
driver/
├── sqlite.go         ← package driver (engine init, registration)
├── conn.go
├── stmt.go
├── rows.go
├── tx.go
├── result.go
├── error.go
├── convert.go
├── vtab.go           ← virtual table support (flat in driver/)
├── backup.go
├── mutex.go          ← renamed to sqlite_mutex.go (avoids clash with Go stdlib)
├── pre_update_hook.go
├── fcntl.go
├── addport.go
├── dmesg.go / nodmesg.go
├── rlimit.go / norlimit.go / rulimit.go
├── driver.go         ← package doc + optional init shim
├── lib/              ← package sqlite3 (auto-generated C→Go, DO NOT rename)
│   └── *.go
├── vfs/              ← package vfs (auto-generated, DO NOT rename)
│   └── *.go
└── vtab/             ← package vtab (vtab helper, DO NOT rename)
    └── *.go
```

> **Sub-package rule exception:** `lib/`, `vfs/`, `vtab/` are auto-generated or
> have mandatory separate package names. They are the only allowed sub-packages.

---

## Steps

### Step 1 — Create `driver/` and move root engine files

```bash
mkdir -p driver

# Move all .go files from root into driver/
# (go.mod, go.sum, README.md stay at root)
mv *.go driver/

# Special rename: mutex.go conflicts with Go stdlib in some compilers
mv driver/mutex.go driver/sqlite_mutex.go
```

### Step 2 — Move sub-packages into `driver/`

```bash
mv lib  driver/lib
mv vfs  driver/vtab
mv vtab driver/vtab
```

Wait — check actual folder names first:
```bash
ls -d */
# Confirm: lib/ vfs/ vtab/ exist at root
```

Then move:
```bash
mv lib  driver/lib
mv vfs  driver/vfs
mv vtab driver/vtab
```

### Step 3 — Rename `package sqlite` → `package driver`

Only rename files directly in `driver/` (not in sub-packages):

```bash
# Replace package declaration in all direct driver/*.go files
sed -i 's|^package sqlite // import "modernc.org/sqlite"|package driver|' driver/*.go
sed -i 's|^package sqlite$|package driver|' driver/*.go
```

Verify no `package sqlite` remains in `driver/*.go`:
```bash
grep "^package sqlite" driver/*.go
# Expected: no output
```

Do NOT rename `driver/lib/`, `driver/vfs/`, `driver/vtab/` — they keep their own package names.

### Step 4 — Update internal import paths in `driver/*.go`

After the move, `lib/`, `vfs/`, `vtab/` are now sub-packages of `driver`:

```bash
# Fix imports in direct driver/ files
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g' driver/*.go
sed -i 's|github.com/cdvelop/sqlite-wasm/vfs|github.com/cdvelop/sqlite-wasm/driver/vfs|g' driver/*.go
sed -i 's|github.com/cdvelop/sqlite-wasm/vtab|github.com/cdvelop/sqlite-wasm/driver/vtab|g' driver/*.go

# Fix imports in vfs/ files (they import lib/)
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g' driver/vfs/*.go

# Fix imports in vtab/ files (they may import lib/)
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g' driver/vtab/*.go
```

Verify no old paths remain:
```bash
grep -r "cdvelop/sqlite-wasm/lib\|cdvelop/sqlite-wasm/vfs\|cdvelop/sqlite-wasm/vtab" \
     driver/ | grep -v "driver/"
# Expected: no output (all paths should be under driver/)
```

### Step 5 — Create `driver/driver.go` (package doc)

Create a minimal package documentation file:

```go
// Package driver embeds the modernc.org/sqlite engine for use as a
// database/sql driver. Its init() function (in sqlite.go) registers
// the "sqlite" driver name automatically on import.
//
// This package is the engine boundary layer. To swap the SQLite engine,
// replace the source files in this package. The adapter (adapter.go in
// the parent module) remains unchanged.
package driver
```

### Step 6 — Build verification

```bash
go build ./...
```

Fix any remaining path or package declaration issues.

### Step 7 — Run tests

```bash
gotest
```

All tests must pass (even if coverage is low at this stage).

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| All engine `.go` files are inside `driver/` (root is empty of `.go` files) | ✅ |
| `driver/*.go` all declare `package driver` | ✅ |
| `driver/lib/*.go` declare `package sqlite3` (unchanged) | ✅ |
| `driver/vfs/*.go` declare `package vfs` (unchanged) | ✅ |
| `driver/vtab/*.go` declare `package vtab` (unchanged) | ✅ |
| All internal imports reference `github.com/cdvelop/sqlite-wasm/driver/...` | ✅ |
| No `modernc.org/sqlite` reference in any `.go` file | ✅ |
| `go build ./...` succeeds | ✅ |
| `gotest` passes | ✅ |
