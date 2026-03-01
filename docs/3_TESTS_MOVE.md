# Phase 3: Move All Tests to `tests/`

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [2_TESTS_PASS.md](2_TESTS_PASS.md) ← must be complete
> **Next:** [4_TESTS_DOMAIN.md](4_TESTS_DOMAIN.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Consolidate all test files into a top-level `tests/` directory as required by the
project rules (>5 test files → move all to `tests/`). Add proper build tags and
a shared setup. All tests must continue to pass.

---

## Rationale

The `driver/` folder is the **deliverable** — it must be clean and contain no test
files. When copied into `tinywasm/sqlite/driver/`, it should contain only
production source. Tests belong in the parent module's `tests/` directory.

---

## Steps

### Step 1 — Create `tests/` directory and move test files

```bash
mkdir -p tests

# Move vfs test (and any others found in driver/)
mv driver/vfs/all_test.go tests/vfs_test.go

# Find and move any other *_test.go files under driver/
find driver/ -name "*_test.go" -exec mv {} tests/ \;
```

### Step 2 — Update package declarations in `tests/`

All files in `tests/` must declare `package tests` (black-box testing convention)
or a specific sub-package. Use `package driver_test` for black-box tests of driver/:

```bash
# Check current package declarations
grep "^package" tests/*.go
```

Rename as appropriate. Black-box tests: `package driver_test`. Internal tests: keep `package driver`.

### Step 3 — Fix import paths in moved test files

After moving, imports like `"github.com/cdvelop/sqlite-wasm/driver/vfs"` must be
explicit since the test is no longer inside `driver/vfs/`:

```bash
# Check what each test file imports and fix paths
grep -h "import" tests/*.go
```

Add explicit import for the package under test:
```go
import (
    "testing"
    "github.com/cdvelop/sqlite-wasm/driver/vfs"
)
```

### Step 4 — Create `tests/setup_test.go`

Shared initialization for all tests in the package:

```go
// Package driver_test contains integration tests for the driver/ package.
package driver_test

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // Setup: ensure tests run from repo root for any relative path dependencies
    os.Exit(m.Run())
}
```

### Step 5 — Run tests

```bash
gotest
```

All tests must pass from their new location.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `driver/` contains zero `*_test.go` files | ✅ |
| All tests are in `tests/` | ✅ |
| `tests/setup_test.go` exists with `TestMain` | ✅ |
| `gotest` passes | ✅ |
