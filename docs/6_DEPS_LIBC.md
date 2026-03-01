# Phase 6: Inline `modernc.org/libc`

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [5_DEPS_SMALL.md](5_DEPS_SMALL.md) ← must be complete
> **Next:** [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Eliminate `modernc.org/libc` from `go.mod` by inlining its source into `driver/libc/`.
This is the largest and most complex step. After this phase, `go.mod` contains
no `modernc.org/*` entries.

---

## Warning

`modernc.org/libc` is a large library (~2000+ files, platform-specific, CGo-free
C standard library emulation in pure Go). Approach incrementally:

1. Copy source files.
2. Fix package declarations and internal imports.
3. Fix compilation errors platform by platform (start with `linux/amd64`).
4. Run `gotest` on each platform before moving to the next.

> **If this phase proves too complex or time-consuming,** accept `modernc.org/libc`
> as a remaining external dep and skip to Phase 7. Document the decision in
> `driver/README.md`. The migration to `tinywasm/sqlite` can still proceed;
> only `libc` remains as an indirect dep.

---

## Steps

### Step 1 — Audit libc usage

```bash
grep -r "modernc.org/libc" driver/ --include="*.go" | grep -v "_test.go"
```

Map which files use libc and what they import:
```bash
grep -r "modernc.org/libc" driver/ --include="*.go" | \
    sed 's|.*"\(modernc.org/libc[^"]*\)".*|\1|' | sort -u
```

### Step 2 — Locate libc source

```bash
GOMOD=$(go env GOMODCACHE)
LIBC_VER=$(grep "modernc.org/libc" go.mod | awk '{print $2}')
SRC="$GOMOD/modernc.org/libc@$LIBC_VER"
echo "Source: $SRC"
ls $SRC/
```

### Step 3 — Plan the sub-package mapping

`modernc.org/libc` exports multiple sub-packages (`libc/sys/types`, `libc/sys/stat`,
etc.). Map each to a `driver/libc/` sub-directory:

| Original | Destination |
|----------|-------------|
| `modernc.org/libc` | `driver/libc/` |
| `modernc.org/libc/sys/types` | `driver/libc/sys/types/` |
| `modernc.org/libc/sys/stat` | `driver/libc/sys/stat/` |
| *(others as needed)* | *(add as compilation requires)* |

### Step 4 — Copy and fix (linux/amd64 first)

```bash
mkdir -p driver/libc
cp $SRC/*.go driver/libc/
cp $SRC/*_linux_amd64.go driver/libc/ 2>/dev/null
rm -f driver/libc/*_test.go

# Copy sub-packages selectively
for subpkg in sys/types sys/stat; do
    mkdir -p driver/libc/$subpkg
    cp $SRC/$subpkg/*.go driver/libc/$subpkg/
    rm -f driver/libc/$subpkg/*_test.go
done
```

### Step 5 — Update import paths

```bash
# In driver/ files: point libc imports to the inlined version
sed -i 's|modernc.org/libc|github.com/cdvelop/sqlite-wasm/driver/libc|g' \
    driver/*.go driver/lib/*.go
```

### Step 6 — Iterative build fix

```bash
# Build only for linux/amd64 first to limit complexity
GOOS=linux GOARCH=amd64 go build ./driver/...
```

Fix compilation errors one by one. Common issues:
- Missing platform files (add them from the source as needed).
- Internal `libc` self-references (apply the same `sed` to `driver/libc/*.go`).
- Circular imports (reorganize sub-packages if needed).

### Step 7 — Expand to other platforms

After `linux/amd64` compiles:
```bash
GOOS=darwin GOARCH=arm64 go build ./driver/...
GOOS=windows GOARCH=amd64 go build ./driver/...
```

Fix similar issues per platform.

### Step 8 — Remove from `go.mod`

```bash
go mod tidy
grep "modernc.org" go.mod
# Expected: no output
```

### Step 9 — Run tests

```bash
gotest
```

Coverage must remain ≥ 90%.

---

## Decision Gate

Before starting this phase, estimate effort:
```bash
GOMOD=$(go env GOMODCACHE)
LIBC_VER=$(grep "modernc.org/libc" go.mod | awk '{print $2}')
find "$GOMOD/modernc.org/libc@$LIBC_VER" -name "*.go" ! -name "*_test.go" | wc -l
```

If > 500 files: consider deferring to a separate dedicated effort and document
`modernc.org/libc` as an accepted remaining dependency.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `modernc.org/libc` absent from `go.mod` | ✅ |
| `driver/libc/` exists with inlined source | ✅ |
| `go build ./...` succeeds on linux/amd64 | ✅ |
| `gotest` passes with ≥ 90% coverage | ✅ |
