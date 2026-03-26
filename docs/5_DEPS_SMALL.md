# Phase 5: Inline Small `modernc.org` Dependencies

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [4_TESTS_DOMAIN.md](4_TESTS_DOMAIN.md) ← must be complete
> **Next:** [6_DEPS_LIBC.md](6_DEPS_LIBC.md)

## Prerequisites

```bash
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

> **`modernc.org/fileutil` is NOT present in `go.mod` — skip it entirely.**
> **Do NOT attempt `modernc.org/libc` in this phase.** It is large and complex — see Phase 6.

---

## Strategy

For each package:
1. Locate its source (`go env GOMODCACHE` → `modernc.org/<pkg>@<version>/`).
2. Copy the needed `.go` files into a **single subdirectory** inside `driver/` (e.g. `driver/mathutil/`). Do NOT nest further.
3. Keep the original package name (e.g. `package mathutil`).
4. Update import paths in `driver/*.go` to point to the local copy.
5. Run `go mod tidy` to confirm the deps are removed from `go.mod`.
6. Run `go test ./...`.

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
SRC="$(go env GOMODCACHE)/modernc.org/mathutil@$PKG_VER"
DEST="driver/mathutil"

mkdir -p $DEST
# Copy non-test .go files into driver/mathutil/
cp $SRC/*.go $DEST/
rm -f $DEST/*_test.go
```

> Repeat the same for `modernc.org/memory` into `driver/memory/`.

### Step 4 — Update import paths in `driver/`

```bash
# Redirect imports to local copies
sed -i 's|modernc.org/mathutil|github.com/cdvelop/sqlite-wasm/driver/mathutil|g' driver/*.go
sed -i 's|modernc.org/memory|github.com/cdvelop/sqlite-wasm/driver/memory|g' driver/*.go
```

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
go test ./...
```

Coverage must remain ≥ 90%.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `modernc.org/mathutil` absent from `go.mod` | ✅ |
| `modernc.org/memory` absent from `go.mod` | ✅ |
| `driver/mathutil/` and `driver/memory/` exist (one level only) | ✅ |
| `go build ./...` succeeds | ✅ |
| `go test ./...` passes with ≥ 90% coverage | ✅ |
