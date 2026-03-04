# Phase 1C: Move `vfs/` and `vtab/` into `driver/`

> **Master Plan:** [REFACTOR_PLAN.md](REFACTOR_PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [1B_MOVE_LIB.md](1B_MOVE_LIB.md) ← must be ✅ committed
> **Next:** [1D_RENAME_PKG.md](1D_RENAME_PKG.md)

## Goal

Move `vfs/` (~24 files, package `vfs`) and `vtab/` (~2 files, package `vtab`)
under `driver/`. Do **not** rename any package or update any import path in this
sub-stage — only move the directories.

---

## Steps

### Step 1 — Move `vfs/` under `driver/`

```bash
git mv vfs driver/vfs
```

### Step 2 — Move `vtab/` under `driver/`

```bash
git mv vtab driver/vtab
```

### Step 3 — Verify

```bash
ls driver/vfs/ | wc -l   # Expected: ~24
ls driver/vtab/ | wc -l  # Expected: ~2

ls vfs  2>/dev/null && echo "FAIL: vfs/ still at root"  || echo "OK"
ls vtab 2>/dev/null && echo "FAIL: vtab/ still at root" || echo "OK"
```

### Step 4 — Commit

```bash
git add -A
git commit -m "move(1c): vfs/ + vtab/ → driver/"
```

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `driver/vfs/` exists with original files | 🔲 |
| `driver/vtab/` exists with original files | 🔲 |
| `vfs/` and `vtab/` no longer exist at repo root | 🔲 |
| No file content changed (package names and imports unchanged) | 🔲 |
| Commit created | 🔲 |
