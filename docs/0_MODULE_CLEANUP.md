# Phase 0: Module Cleanup

> **Master Plan:** [PLAN.md](PLAN.md)
> **Next:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)

## Prerequisites

```bash
# Install gotest (required — do NOT use go test directly)
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Fix the inherited `go.mod` module name and resolve the package name conflict that
prevents `go build ./...` from compiling. No files are moved in this phase.

---

## Problem

The forked source was copied as-is. It has two issues:

1. **Wrong module name:** `go.mod` still declares `module modernc.org/sqlite`.
   All internal imports (e.g. `"modernc.org/sqlite/lib"`) reference this name.

2. **Package conflict:** `sqlite-wasm.go` declares `package sqlitewasm` while
   all other root files declare `package sqlite`. Go does not allow two package
   names in the same directory → `go build` fails.

---

## Steps

### Step 1 — Fix `go.mod` module name

```bash
# In repo root
sed -i 's|^module modernc.org/sqlite|module github.com/cdvelop/sqlite-wasm|' go.mod
```

Verify:
```bash
head -1 go.mod
# Expected: module github.com/cdvelop/sqlite-wasm
```

### Step 2 — Remove the stub file

`sqlite-wasm.go` (`package sqlitewasm`) conflicts with `package sqlite` in the
same directory. Delete it:

```bash
rm sqlite-wasm.go
```

### Step 3 — Update internal import paths (root-level files)

All `.go` files at root reference themselves as `modernc.org/sqlite/lib`,
`modernc.org/sqlite/vfs`, `modernc.org/sqlite/vtab`. Replace with the new module path:

```bash
# Replace all self-referencing modernc.org/sqlite/* imports
find . -maxdepth 1 -name "*.go" \
  -exec sed -i 's|modernc.org/sqlite/lib|github.com/cdvelop/sqlite-wasm/lib|g' {} \; \
  -exec sed -i 's|modernc.org/sqlite/vfs|github.com/cdvelop/sqlite-wasm/vfs|g' {} \; \
  -exec sed -i 's|modernc.org/sqlite/vtab|github.com/cdvelop/sqlite-wasm/vtab|g' {} \;

# Also fix vfs/ sub-package (it imports modernc.org/sqlite/lib)
find vfs/ -name "*.go" \
  -exec sed -i 's|modernc.org/sqlite/lib|github.com/cdvelop/sqlite-wasm/lib|g' {} \;

# Also fix vtab/ sub-package
find vtab/ -name "*.go" \
  -exec sed -i 's|modernc.org/sqlite/lib|github.com/cdvelop/sqlite-wasm/lib|g' {} \;
```

Verify no `modernc.org/sqlite` self-references remain (only `modernc.org/libc` and
other deps are acceptable):

```bash
grep -r "modernc.org/sqlite" --include="*.go" .
# Should show NO results
```

### Step 4 — Run `go mod tidy`

```bash
go mod tidy
```

### Step 5 — Verify build

```bash
go build ./...
```

Fix any remaining import or package issues before proceeding.

### Step 6 — Run tests

```bash
gotest
```

Tests that existed before must still pass.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `go.mod` first line: `module github.com/cdvelop/sqlite-wasm` | ✅ |
| `sqlite-wasm.go` deleted | ✅ |
| No `modernc.org/sqlite` self-references in any `.go` file | ✅ |
| `go build ./...` succeeds | ✅ |
| `gotest` passes | ✅ |
