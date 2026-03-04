# Phase 1: Organize All Engine Source into `driver/`

> **Master Plan:** [REFACTOR_PLAN.md](REFACTOR_PLAN.md)
> **Previous:** [0_MODULE_CLEANUP.md](0_MODULE_CLEANUP.md) ← must be ✅ complete
> **Next:** [2_TESTS_PASS.md](2_TESTS_PASS.md)

## Goal

Move all SQLite engine source files from the repo root into `driver/`, rename
`package sqlite` → `package driver`, and update all internal import paths.

**This phase is split into 6 atomic sub-stages. Each sub-stage must be committed
before starting the next.** This prevents the session diff buffer from overflowing
on large moves (70+ files).

---

## Sub-stages

| File | Goal | Commit scope |
|------|------|--------------|
| [1A_MOVE_ROOT.md](1A_MOVE_ROOT.md) | `git mv *.go driver/` + rename `mutex.go` | ~20 files |
| [1B_MOVE_LIB.md](1B_MOVE_LIB.md) | `git mv lib/ driver/lib/` | ~48 files |
| [1C_MOVE_VFS_VTAB.md](1C_MOVE_VFS_VTAB.md) | `git mv vfs/ driver/vfs/` + `vtab/` | ~26 files |
| [1D_RENAME_PKG.md](1D_RENAME_PKG.md) | `sed` `package sqlite` → `package driver` in `driver/*.go` | ~20 files (text) |
| [1E_FIX_IMPORTS.md](1E_FIX_IMPORTS.md) | Fix all internal import paths under `driver/` | ~46 files (text) |
| [1F_BUILD_VERIFY.md](1F_BUILD_VERIFY.md) | Create `driver/driver.go` + `go build ./...` + `gotest` | 1 new file |

---

## Phase 1 Acceptance Criteria (all sub-stages done)

| Criterion | Check |
|-----------|-------|
| All engine `.go` files are inside `driver/` (root has no `.go` files) | 🔲 |
| `driver/*.go` all declare `package driver` | 🔲 |
| `driver/lib/*.go` declare `package sqlite3` (unchanged) | 🔲 |
| `driver/vfs/*.go` declare `package vfs` (unchanged) | 🔲 |
| `driver/vtab/*.go` declare `package vtab` (unchanged) | 🔲 |
| All internal imports reference `github.com/cdvelop/sqlite-wasm/driver/...` | 🔲 |
| `go build ./...` succeeds | 🔲 |
| `gotest` passes | 🔲 |
