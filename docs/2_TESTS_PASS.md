# Phase 2: Tests Pass with New `driver/` Structure

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md) ← must be complete
> **Next:** [3_TESTS_MOVE.md](3_TESTS_MOVE.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Ensure all tests that existed in the original fork pass correctly with the new
`driver/` package structure. This phase does **not** add new tests — it only
fixes failures caused by the reorganization.

---

## Context

After Phase 1, the test file `driver/vfs/all_test.go` is the main candidate.
It was originally at `vfs/all_test.go` and imports the `vfs` package.

Expected issues to fix:
1. `vfs/all_test.go` import paths may reference old module path.
2. Any test helper that assumed files were at the repo root.
3. Build tag mismatches after package rename.

---

## Steps

### Step 1 — Run `gotest` and capture failures

```bash
gotest 2>&1 | tee /tmp/phase2_failures.txt
```

### Step 2 — Fix import paths in test files

```bash
# Fix vfs test imports
sed -i 's|github.com/cdvelop/sqlite-wasm/lib|github.com/cdvelop/sqlite-wasm/driver/lib|g' \
    driver/vfs/all_test.go
```

Repeat for any other test file that has import path errors.

### Step 3 — Fix any `testdata/` references

Some tests reference `testdata/` directories by relative path. After the move,
paths like `"../testdata"` may be broken.

Check:
```bash
grep -r "testdata\|\.\./" driver/ --include="*_test.go"
```

If found: adjust paths or add a `TestMain` that `os.Chdir` to the correct directory.

### Step 4 — Fix build tags if needed

```bash
grep -r "//go:build\|// +build" driver/ --include="*_test.go"
```

Ensure build tags are compatible with the new package name (`driver` instead of `sqlite`).

### Step 5 — Run full test suite

```bash
gotest
```

All inherited tests must pass. Coverage percentage is not a target for this phase.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `gotest` exits with code 0 | ✅ |
| No test failures introduced by the reorganization | ✅ |
| No test files removed (only fixed) | ✅ |
