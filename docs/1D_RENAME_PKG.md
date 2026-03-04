# Phase 1D: Rename `package sqlite` → `package driver`

> **Master Plan:** [PLAN.md](PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [1C_MOVE_VFS_VTAB.md](1C_MOVE_VFS_VTAB.md) ← must be ✅ committed
> **Next:** [1E_FIX_IMPORTS.md](1E_FIX_IMPORTS.md)

## Pre-condition Check

> **IMPORTANT:** The **module name** (`github.com/cdvelop/sqlite-wasm` in `go.mod`) and
> the **package name** (`package sqlite` in `.go` files) are two different things.
> This phase changes **only the package declarations** — it does **not** touch `go.mod`.

Verify the previous phases are committed and `driver/*.go` still declare `package sqlite`:

```bash
# Must show files in driver/ (not at root)
ls driver/*.go | head -5

# Must print lines with 'package sqlite' — if empty, 1A-1C were not committed
grep -l "^package sqlite" driver/*.go
# Expected: several files listed
```

If `grep` returns no files, **STOP** — the previous phases are incomplete.

---


Rename the package declaration in all `driver/*.go` files (not sub-packages)
from `package sqlite` to `package driver`. No file moves, no import changes.

---

## Steps

### Step 1 — Replace package declarations

```bash
# Handles both forms present in the source
sed -i 's|^package sqlite // import "modernc.org/sqlite"|package driver|' driver/*.go
sed -i 's|^package sqlite$|package driver|' driver/*.go
```

### Step 2 — Verify no `package sqlite` remains in `driver/*.go`

```bash
grep "^package sqlite" driver/*.go
# Expected: no output
```

### Step 3 — Verify sub-packages are untouched

```bash
grep "^package" driver/lib/*.go  | sort -u
# Expected: package sqlite3

grep "^package" driver/vfs/*.go  | sort -u
# Expected: package vfs

grep "^package" driver/vtab/*.go | sort -u
# Expected: package vtab
```

### Step 4 — Commit

```bash
git add driver/*.go
git commit -m "rename(1d): package sqlite → package driver in driver/*.go"
```

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `grep "^package sqlite" driver/*.go` returns no output | 🔲 |
| All `driver/*.go` declare `package driver` | 🔲 |
| `driver/lib/`, `driver/vfs/`, `driver/vtab/` package names unchanged | 🔲 |
| Commit created | 🔲 |
