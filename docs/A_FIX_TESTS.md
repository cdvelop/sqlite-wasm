# Phase A: Fix `TestBackup` Deadlock

> **Master Plan:** [PLAN.md](PLAN.md)
> **Next:** [8_MIGRATION.md](8_MIGRATION.md)

## Problem

`go test ./...` panics with a deadlock:

```
panic: test timed out after 1m0s
    running tests: TestBackup
```

Root cause: `tests/backup_test.go:28` calls `src.Conn(nil)`.
`database/sql.DB.Conn` requires a non-nil `context.Context`.
Passing `nil` causes `db.conn()` to panic inside `database/sql`, which
triggers a deferred `db.Close()` that deadlocks on the same mutex.

## Fix

In `tests/backup_test.go`, replace:

```go
conn, err := src.Conn(nil)
```

with:

```go
conn, err := src.Conn(context.Background())
```

And add `"context"` to the import block.

## Steps

### Step 1 — Apply the fix

Edit `tests/backup_test.go`:
- Add `"context"` import.
- Change `src.Conn(nil)` → `src.Conn(context.Background())`.

### Step 2 — Run tests

```bash
go test ./...
```

All tests must pass. Coverage ≥ 90%.

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `TestBackup` passes without panic or timeout | ✅ |
| `go test ./...` exits 0 | ✅ |
