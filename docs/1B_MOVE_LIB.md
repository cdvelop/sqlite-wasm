# Phase 1B: Move `lib/` into `driver/lib/`

> **Master Plan:** [REFACTOR_PLAN.md](REFACTOR_PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [1A_MOVE_ROOT.md](1A_MOVE_ROOT.md) ← must be ✅ committed
> **Next:** [1C_MOVE_VFS_VTAB.md](1C_MOVE_VFS_VTAB.md)

## Goal

Move the auto-generated `lib/` directory (package `sqlite3`, ~48 files) under
`driver/`. Do **not** rename any package or update any import path in this
sub-stage — only move the directory.

---

## Steps

### Step 1 — Move `lib/` under `driver/`

```bash
git mv lib driver/lib
```

### Step 2 — Verify

```bash
ls driver/lib/ | head -5
# Expected: defs.go, hooks.go, sqlite_linux_amd64.go, etc.

ls lib 2>/dev/null && echo "FAIL: lib/ still at root" || echo "OK: lib/ moved"
```

### Step 3 — Commit

```bash
git add -A
git commit -m "move(1b): lib/ → driver/lib/"
```

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `driver/lib/` exists and contains the original files | 🔲 |
| `lib/` no longer exists at repo root | 🔲 |
| No file content changed (package names and imports unchanged) | 🔲 |
| Commit created | 🔲 |
