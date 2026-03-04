# Phase 1D: Rename `package sqlite` → `package driver`

> **Master Plan:** [REFACTOR_PLAN.md](REFACTOR_PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [1C_MOVE_VFS_VTAB.md](1C_MOVE_VFS_VTAB.md) ← must be ✅ committed
> **Next:** [1E_FIX_IMPORTS.md](1E_FIX_IMPORTS.md)

## Goal

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
