# Phase 8: Migration Script & Final Validation

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md) ← must be complete

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Write and validate `scripts/migrate_to_tinywasm.sh` — the script that copies the
`driver/` folder into `tinywasm/sqlite/driver/` and substitutes the import paths.
Run a dry-run to confirm the migration works end-to-end before opening the PR.

After this phase: open the PR from `cdvelop/sqlite-wasm` → `tinywasm/sqlite`.

---

## Migration Script

Create `scripts/migrate_to_tinywasm.sh`:

```bash
#!/usr/bin/env bash
# migrate_to_tinywasm.sh
# Copies driver/ from cdvelop/sqlite-wasm into the tinywasm/sqlite repo
# and substitutes import paths for the new location.
#
# Usage:
#   ./scripts/migrate_to_tinywasm.sh <path/to/tinywasm/sqlite>
#
# Example:
#   ./scripts/migrate_to_tinywasm.sh /home/cesar/Dev/Project/tinywasm/sqlite

set -euo pipefail

SRC_DRIVER="$(cd "$(dirname "$0")/.." && pwd)/driver"
DEST_REPO="${1:?Usage: $0 <path/to/tinywasm/sqlite>}"
DEST_DRIVER="$DEST_REPO/driver"

OLD_PATH="github.com/cdvelop/sqlite-wasm/driver"
NEW_PATH="github.com/tinywasm/sqlite/driver"

echo "==> Source:      $SRC_DRIVER"
echo "==> Destination: $DEST_DRIVER"
echo "==> Old path:    $OLD_PATH"
echo "==> New path:    $NEW_PATH"
echo ""

# Step 1: Clean destination driver/
if [ -d "$DEST_DRIVER" ]; then
    echo "==> Removing existing $DEST_DRIVER ..."
    rm -rf "$DEST_DRIVER"
fi

# Step 2: Copy driver/ verbatim
echo "==> Copying driver/ ..."
cp -r "$SRC_DRIVER" "$DEST_DRIVER"

# Step 3: Substitute import paths in all .go files
echo "==> Substituting import paths ..."
find "$DEST_DRIVER" -name "*.go" \
    -exec sed -i "s|$OLD_PATH|$NEW_PATH|g" {} \;

# Step 4: Verify no old paths remain
echo "==> Verifying substitution ..."
REMAINING=$(grep -r "$OLD_PATH" "$DEST_DRIVER" --include="*.go" || true)
if [ -n "$REMAINING" ]; then
    echo "ERROR: Old import paths still found:"
    echo "$REMAINING"
    exit 1
fi

echo ""
echo "==> Migration complete."
echo "==> Next steps:"
echo "    1. cd $DEST_REPO"
echo "    2. go mod tidy"
echo "    3. go build ./..."
echo "    4. gotest"
echo "    5. Open a PR from cdvelop/sqlite-wasm to tinywasm/sqlite"
```

Make executable:
```bash
chmod +x scripts/migrate_to_tinywasm.sh
```

---

## Steps

### Step 1 — Write the migration script

Create `scripts/migrate_to_tinywasm.sh` as above.

### Step 2 — Dry-run into a temp directory

Validate the script without touching the real `tinywasm/sqlite`:

```bash
mkdir -p /tmp/sqlite_migration_test
./scripts/migrate_to_tinywasm.sh /tmp/sqlite_migration_test

# Verify output
ls /tmp/sqlite_migration_test/driver/
grep -r "cdvelop/sqlite-wasm" /tmp/sqlite_migration_test/driver/ --include="*.go"
# Expected: no output (all paths replaced)
```

### Step 3 — Write `driver/README.md`

Document the package for the next maintainer:

```markdown
# driver — SQLite Engine Package

This package is a self-contained embedding of the SQLite engine
(originally `modernc.org/sqlite`) for use as a `database/sql` driver.

## Registration

The `init()` function in `sqlite.go` registers the `"sqlite"` driver name
automatically on blank import:

    import _ "github.com/tinywasm/sqlite/driver"

## Sub-packages

| Package | Description |
|---------|-------------|
| `driver/lib/` | `package sqlite3` — auto-generated C→Go SQLite amalgamation |
| `driver/vfs/` | `package vfs` — virtual file system layer |
| `driver/vtab/` | `package vtab` — virtual table helpers |

## Dependencies

After full refinement, only `golang.org/x/sys` may remain as an indirect
low-level OS dep (required by the embedded libc layer).

## Origin

Built and refined in `github.com/cdvelop/sqlite-wasm`.
```

### Step 4 — Final `gotest` in this repo

```bash
gotest
```

All tests pass, coverage ≥ 90%.

### Step 5 — Run actual migration into `tinywasm/sqlite`

```bash
TINYWASM_SQLITE=/home/cesar/Dev/Project/tinywasm/sqlite
./scripts/migrate_to_tinywasm.sh $TINYWASM_SQLITE
```

Then in `tinywasm/sqlite`:
```bash
cd $TINYWASM_SQLITE
go mod tidy
go build ./...
gotest
```

### Step 6 — Verify in `tinywasm/sqlite`

Run `gotest` in the target repo to confirm the migration works end-to-end.
Jules will handle committing and pushing.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `scripts/migrate_to_tinywasm.sh` exists and is executable | ✅ |
| Dry-run produces no `cdvelop/sqlite-wasm` references in output | ✅ |
| `driver/README.md` documents the package and sub-packages | ✅ |
| After migration, `gotest` passes in `tinywasm/sqlite` | ✅ |
| PR from `cdvelop/sqlite-wasm` → `tinywasm/sqlite` is open | ✅ |
