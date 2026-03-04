# Phase 1F: Create `driver/driver.go` + Build Verification

> **Master Plan:** [PLAN.md](PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [1E_FIX_IMPORTS.md](1E_FIX_IMPORTS.md) ← must be ✅ committed
> **Next:** [2_TESTS_PASS.md](2_TESTS_PASS.md)

## Goal

Create the `driver/driver.go` package documentation file, confirm the build
compiles, and verify all tests pass. This is the final gate for Phase 1.

---

## Steps

### Step 1 — Create `driver/driver.go`

Create the file with exactly this content:

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

### Step 2 — Build

```bash
go build ./...
```

Fix any remaining import or package declaration issues before proceeding.

### Step 3 — Run tests

```bash
gotest
```

All tests must pass.

### Step 4 — Commit

```bash
git add driver/driver.go
git commit -m "add(1f): driver/driver.go; go build ./... and gotest pass"
```

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `driver/driver.go` exists with `package driver` | 🔲 |
| `go build ./...` exits with code 0 | 🔲 |
| `gotest` passes with no failures | 🔲 |
| Commit created | 🔲 |

---

## ⛔ DISPATCH D-3 COMPLETE — STOP HERE

This is the **last phase of Dispatch D-3**. Your task is complete.

**Do NOT open, read, or execute `2_TESTS_PASS.md` or any subsequent file.**

Report the following to the user and wait for further instructions:

```
✅ D-3 done. Commit created:
  - add(1f): driver/driver.go; go build ./... and gotest pass

gotest output:
  <paste full gotest output here>

Awaiting authorization for Dispatch D-4 (2_TESTS_PASS.md).
```
