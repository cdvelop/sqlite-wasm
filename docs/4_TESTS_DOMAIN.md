# Phase 4: Subdivide Tests by Domain

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [3_TESTS_MOVE.md](3_TESTS_MOVE.md) ← must be complete
> **Next:** [5_DEPS_SMALL.md](5_DEPS_SMALL.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Split tests into domain-focused files. Write new tests to reach **≥ 90% coverage**
of `driver/`. Each test file covers one area of the SQLite driver API.

---

## Domain Split

| Test File | Domain | Key Functions Covered |
|-----------|--------|-----------------------|
| `tests/setup_test.go` | Shared | `TestMain`, helpers, DB open/close |
| `tests/conn_test.go` | Connections | `Open`, `Close`, `Ping`, connection parameters |
| `tests/exec_test.go` | Execution | `Exec`, `Query`, DDL (CREATE TABLE, DROP) |
| `tests/stmt_test.go` | Statements | `Prepare`, `Stmt.Exec`, `Stmt.Query`, `Stmt.Close` |
| `tests/tx_test.go` | Transactions | `Begin`, `Commit`, `Rollback`, nested tx |
| `tests/rows_test.go` | Row scanning | `Rows.Next`, `Rows.Scan`, `Rows.Close`, column types |
| `tests/backup_test.go` | Backup API | `Backup`, backup to `:memory:` |
| `tests/vtab_test.go` | Virtual tables | Virtual table registration, query |
| `tests/vfs_test.go` | VFS | VFS registration, custom file system ops |
| `tests/error_test.go` | Error handling | Constraint violations, syntax errors, closed DB |

---

## Steps

### Step 1 — Audit current coverage

```bash
gotest -cover ./driver/...
```

Identify which files have low coverage. Prioritize those.

### Step 2 — Write `tests/conn_test.go`

```go
//go:build !wasm

package driver_test

import (
    "database/sql"
    "testing"
    _ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestOpen_InMemory(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    if err := db.Ping(); err != nil {
        t.Fatal(err)
    }
}

func TestOpen_InvalidPath(t *testing.T) {
    db, err := sql.Open("sqlite", "/nonexistent/path/db.sqlite")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    // Ping should fail for invalid path
    if err := db.Ping(); err == nil {
        t.Fatal("expected error for invalid path")
    }
}
```

### Step 3 — Write `tests/exec_test.go`

```go
//go:build !wasm

package driver_test

import (
    "database/sql"
    "testing"
    _ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestExec_CreateAndInsert(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil { t.Fatal(err) }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)`)
    if err != nil { t.Fatal(err) }

    _, err = db.Exec(`INSERT INTO users (name) VALUES (?)`, "alice")
    if err != nil { t.Fatal(err) }

    var name string
    err = db.QueryRow(`SELECT name FROM users WHERE id=1`).Scan(&name)
    if err != nil { t.Fatal(err) }
    if name != "alice" {
        t.Fatalf("got %q, want %q", name, "alice")
    }
}
```

### Step 4 — Write remaining domain test files

Follow the same pattern for `stmt_test.go`, `tx_test.go`, `rows_test.go`,
`backup_test.go`, `error_test.go`. See domain table above for coverage targets.

### Step 5 — Check coverage

```bash
gotest -cover ./driver/...
# Target: ≥ 90%
```

If below 90%, identify uncovered functions and add targeted tests.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| Each domain has a dedicated test file in `tests/` | ✅ |
| All test files have `//go:build !wasm` tag | ✅ |
| `gotest -cover ./driver/...` reports ≥ 90% | ✅ |
| `gotest` passes with no failures | ✅ |
