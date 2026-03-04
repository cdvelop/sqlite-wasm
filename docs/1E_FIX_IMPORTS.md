# Phase 1E: Fix Internal Import Paths

> **Master Plan:** [PLAN.md](PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [1D_RENAME_PKG.md](1D_RENAME_PKG.md) ← must be ✅ committed
> **Next:** [1F_BUILD_VERIFY.md](1F_BUILD_VERIFY.md)

## Goal

Update all internal import paths to reflect the new location of `lib/`, `vfs/`,
and `vtab/` under `driver/`. No package renames, no file moves.

---

## Steps

### Step 1 — Fix imports in `driver/*.go`

```bash
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g'   driver/*.go
sed -i 's|github.com/cdvelop/sqlite-wasm/vfs|github.com/cdvelop/sqlite-wasm/driver/vfs|g'   driver/*.go
sed -i 's|github.com/cdvelop/sqlite-wasm/vtab|github.com/cdvelop/sqlite-wasm/driver/vtab|g' driver/*.go
```

### Step 2 — Fix imports in `driver/vfs/*.go` (imports `lib/`)

```bash
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g' driver/vfs/*.go
```

### Step 3 — Fix imports in `driver/vtab/*.go` (may import `lib/`)

```bash
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g' driver/vtab/*.go
```

### Step 4 — Verify no old paths remain

```bash
grep -r "cdvelop/sqlite-wasm/lib\|cdvelop/sqlite-wasm/vfs\|cdvelop/sqlite-wasm/vtab" driver/ \
  | grep -v "driver/lib\|driver/vfs\|driver/vtab"
# Expected: no output
```

### Step 5 — Commit

```bash
git add driver/
git commit -m "fix(1e): update internal import paths to driver/lib, driver/vfs, driver/vtab"
```

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| No `cdvelop/sqlite-wasm/lib` (without `driver/`) in any `.go` file | 🔲 |
| No `cdvelop/sqlite-wasm/vfs` (without `driver/`) in any `.go` file | 🔲 |
| No `cdvelop/sqlite-wasm/vtab` (without `driver/`) in any `.go` file | 🔲 |
| Commit created | 🔲 |

---

## ⛔ DISPATCH D-2 COMPLETE — STOP HERE

This is the **last phase of Dispatch D-2**. Your task is complete.

**Do NOT open, read, or execute `1F_BUILD_VERIFY.md` or any subsequent file.**

Report the following to the user and wait for further instructions:

```
✅ D-2 done. Commits created:
  - rename(1d): package sqlite → package driver in driver/*.go
  - fix(1e): update internal import paths to driver/lib, driver/vfs, driver/vtab

Verification output from Step 4 grep:
  <paste grep output here — expected: empty / no matches>

Awaiting authorization for Dispatch D-3 (1F_BUILD_VERIFY.md).
```
