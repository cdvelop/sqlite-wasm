# Phase 5: Inline Small `modernc.org` Dependencies

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [4_TESTS_DOMAIN.md](4_TESTS_DOMAIN.md) ← must be complete
> **Next:** [6_DEPS_LIBC.md](6_DEPS_LIBC.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Remove the smaller `modernc.org/*` transitive dependencies from `go.mod` by inlining
their source directly into `driver/`. This reduces the external dependency footprint
before tackling the larger `modernc.org/libc` in Phase 6.

---

## Target Dependencies

| Package | Size | Role |
|---------|------|------|
| `modernc.org/mathutil` | Small | Math helpers (used by libc/sqlite) |
| `modernc.org/memory` | Small | Memory management primitives |
| `modernc.org/fileutil` | Small | File utility helpers |

> **Do NOT attempt `modernc.org/libc` in this phase.** It is large and complex — see Phase 6.

---

## Strategy

For each package:
1. Locate its source (run `go env GOPATH` → `pkg/mod/modernc.org/<pkg>@<version>/`).
2. Copy the needed `.go` files into a `driver/<pkg>/` sub-directory.
3. **Rename the package declaration** to match the sub-directory name.
4. Update import paths in `driver/*.go` files that used the original module.
5. Run `go mod tidy` to confirm the dependency is removed.
6. Run `gotest`.

---

## Steps

### Step 1 — Audit which files use these deps

```bash
grep -r "modernc.org/mathutil\|modernc.org/memory\|modernc.org/fileutil" \
     driver/ --include="*.go" | grep -v "lib/"
```

Note which files import each package. These are the files that need import updates.

### Step 2 — Locate source in module cache

```bash
GOPATH=$(go env GOPATH)
GOMOD=$(go env GOMODCACHE)

# Find exact versions from go.mod
grep "modernc.org/mathutil\|modernc.org/memory\|modernc.org/fileutil" go.mod
```

### Step 3 — Copy & inline each package

For each package (example with `mathutil`):

```bash
PKG_VER="v1.7.1"  # adjust from go.mod
SRC="$GOMOD/modernc.org/mathutil@$PKG_VER"
DEST="driver/mathutil"

mkdir -p $DEST
# Copy only non-test .go files
cp $SRC/*.go $DEST/
rm -f $DEST/*_test.go

# Rename package declaration
sed -i 's|^package mathutil|package mathutil|' $DEST/*.go
# (no rename needed if package name matches dir name)
```

### Step 4 — Update import paths in `driver/`

```bash
# Example for mathutil
sed -i 's|modernc.org/mathutil|github.com/cdvelop/sqlite-wasm/driver/mathutil|g' driver/*.go
```

Repeat for each inlined package.

### Step 5 — Remove from `go.mod`

```bash
go mod tidy
```

Verify the packages are gone from `go.mod` and `go.sum`:
```bash
grep "modernc.org/mathutil\|modernc.org/memory\|modernc.org/fileutil" go.mod go.sum
# Expected: no output
```

### Step 6 — Build and test

```bash
go build ./...
gotest
```

Coverage must remain ≥ 90%.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `modernc.org/mathutil` absent from `go.mod` | ✅ |
| `modernc.org/memory` absent from `go.mod` | ✅ |
| `modernc.org/fileutil` absent from `go.mod` | ✅ |
| `driver/mathutil/`, `driver/memory/`, `driver/fileutil/` exist | ✅ |
| `go build ./...` succeeds | ✅ |
| `gotest` passes with ≥ 90% coverage | ✅ |
