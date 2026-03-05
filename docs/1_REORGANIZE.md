# Phase 1: Reorganize Engine Source → `driver/`

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [0_MODULE_CLEANUP.md](0_MODULE_CLEANUP.md) ← must be ✅ complete
> **Next:** [2_TESTS_PASS.md](2_TESTS_PASS.md)

## Goal

Move **all** engine source files into `driver/`, rename `package sqlite` → `package driver`
in the root-level files, and fix all internal import paths — in a **single atomic commit**.

This entire phase is executed by running one bash script. Do NOT split it into smaller steps.

---

## Prerequisites

Confirm Phase 0 is complete and the module builds from the current state:

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
go build ./...
```

If `go build ./...` fails, **STOP** — Phase 0 is not done.

---

## Step 1 — Create the reorganization script

Create the file `scripts/reorganize.sh` with exactly this content:

```bash
#!/usr/bin/env bash
# reorganize.sh — moves modernc engine source into driver/ package
# Safe to re-run: checks for existing structure before acting.
set -euo pipefail

MODULE="github.com/cdvelop/sqlite-wasm"

echo "=== Phase 1: Reorganize engine source into driver/ ==="

# ── 1A: Move root *.go files ─────────────────────────────────────────────────
if ls driver/*.go &>/dev/null 2>&1; then
  echo "[1A] SKIP: driver/*.go already exist"
else
  echo "[1A] Moving root *.go → driver/"
  mkdir -p driver
  git mv *.go driver/
  # Rename mutex.go to avoid stdlib name clash
  [ -f driver/mutex.go ] && git mv driver/mutex.go driver/sqlite_mutex.go
  echo "[1A] OK"
fi

# ── 1B: Move lib/ ────────────────────────────────────────────────────────────
if [ -d driver/lib ]; then
  echo "[1B] SKIP: driver/lib/ already exists"
else
  echo "[1B] Moving lib/ → driver/lib/"
  git mv lib driver/lib
  echo "[1B] OK"
fi

# ── 1C: Move vfs/ and vtab/ ──────────────────────────────────────────────────
if [ -d driver/vfs ]; then
  echo "[1C] SKIP: driver/vfs/ already exists"
else
  echo "[1C] Moving vfs/ → driver/vfs/"
  git mv vfs driver/vfs
  echo "[1C] OK"
fi

if [ -d driver/vtab ]; then
  echo "[1C] SKIP: driver/vtab/ already exists"
else
  echo "[1C] Moving vtab/ → driver/vtab/"
  git mv vtab driver/vtab
  echo "[1C] OK"
fi

# ── 1D: Rename package sqlite → package driver in driver/*.go ────────────────
echo "[1D] Renaming package declarations in driver/*.go"
sed -i 's|^package sqlite // import "modernc.org/sqlite"|package driver|' driver/*.go
sed -i 's|^package sqlite$|package driver|' driver/*.go
PKG_REMAINING=$(grep -l "^package sqlite" driver/*.go 2>/dev/null || true)
if [ -n "$PKG_REMAINING" ]; then
  echo "ERROR [1D]: package sqlite still present in: $PKG_REMAINING"
  exit 1
fi
echo "[1D] OK: all driver/*.go now declare package driver"

# ── 1E: Fix internal import paths ────────────────────────────────────────────
echo "[1E] Fixing internal import paths"
# driver/*.go: lib, vfs, vtab
sed -i "s|${MODULE}/lib|${MODULE}/driver/lib|g"   driver/*.go
sed -i "s|${MODULE}/vfs|${MODULE}/driver/vfs|g"   driver/*.go
sed -i "s|${MODULE}/vtab|${MODULE}/driver/vtab|g" driver/*.go
# driver/vfs/*.go: imports lib
sed -i "s|${MODULE}/lib|${MODULE}/driver/lib|g" driver/vfs/*.go
# driver/vtab/*.go: may import lib
if ls driver/vtab/*.go &>/dev/null 2>&1; then
  sed -i "s|${MODULE}/lib|${MODULE}/driver/lib|g" driver/vtab/*.go
fi
# Verify no old paths remain
STALE=$(grep -r "${MODULE}/lib\|${MODULE}/vfs\|${MODULE}/vtab" driver/ \
  | grep -v "driver/lib\|driver/vfs\|driver/vtab" || true)
if [ -n "$STALE" ]; then
  echo "ERROR [1E]: stale import paths found:"
  echo "$STALE"
  exit 1
fi
echo "[1E] OK: all import paths updated"

# ── 1F: Create driver/driver.go ───────────────────────────────────────────────
if [ ! -f driver/driver.go ]; then
  echo "[1F] Creating driver/driver.go"
  cat > driver/driver.go <<'GOFILE'
// Package driver embeds the modernc.org/sqlite engine for use as a
// database/sql driver. Its init() function (in sqlite.go) registers
// the "sqlite" driver name automatically on import.
//
// This package is the engine boundary layer. To swap the SQLite engine,
// replace the source files in this package. The adapter (adapter.go in
// the parent module) remains unchanged.
package driver
GOFILE
else
  echo "[1F] SKIP: driver/driver.go already exists"
fi

# ── Verify the build ──────────────────────────────────────────────────────────
echo "[BUILD] Running go build ./..."
go build ./...
echo "[BUILD] OK"

# ── Single atomic commit ──────────────────────────────────────────────────────
echo "[COMMIT] Staging all changes"
git add -A
git status --short

echo "[COMMIT] Creating commit"
git commit -m "reorganize(1): engine source → driver/; package driver; fix imports"

echo ""
echo "=== Phase 1 COMPLETE ==="
echo ""
echo "driver/*.go   → $(ls driver/*.go 2>/dev/null | wc -l) files"
echo "driver/lib/   → $(ls driver/lib/ 2>/dev/null | wc -l) files"
echo "driver/vfs/   → $(ls driver/vfs/ 2>/dev/null | wc -l) files"
echo "driver/vtab/  → $(ls driver/vtab/ 2>/dev/null | wc -l) files"
echo ""
echo "Awaiting authorization for Dispatch D-2 (2_TESTS_PASS.md)."
```

---

## Step 2 — Run the script

```bash
chmod +x scripts/reorganize.sh
bash scripts/reorganize.sh
```

Paste the **full output** in your report. The script exits with a non-zero code on any error.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| No `.go` files remain at repo root (only `go.mod`, `go.sum`, `README.md`) | 🔲 |
| `driver/*.go` declare `package driver` (not `package sqlite`) | 🔲 |
| `driver/lib/`, `driver/vfs/`, `driver/vtab/` package names unchanged | 🔲 |
| `grep -r "cdvelop/sqlite-wasm/lib" driver/` returns only paths containing `driver/lib` | 🔲 |
| `go build ./...` exits with code 0 | 🔲 |
| Exactly **one** commit created: `reorganize(1): engine source → driver/; package driver; fix imports` | 🔲 |

---

## ⛔ DISPATCH D-1 COMPLETE — STOP HERE

This is the **last phase of Dispatch D-1**. Your task is complete.

**Do NOT open, read, or execute `2_TESTS_PASS.md` or any subsequent file.**

Report the following to the user and wait for further instructions:

```
✅ D-1 done. Commit: reorganize(1): engine source → driver/; package driver; fix imports

Script output:
  <paste full script output here>

File counts:
  driver/*.go   → N files
  driver/lib/   → N files
  driver/vfs/   → N files
  driver/vtab/  → N files

Awaiting authorization for Dispatch D-2 (2_TESTS_PASS.md).
```
