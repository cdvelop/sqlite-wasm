# PLAN: Remaining Work — Test Fix + Migration

## Purpose

Finalize the `driver/` package and migrate it to `tinywasm/sqlite/driver/`.
`go.mod` is already clean: `modernc.org/libc` is the only direct dependency.
Two tasks remain: fix a failing test, then run the migration.

> **IMPORTANT FOR THE AGENT:** Execute phases in order. After each phase run
> `go test ./...` to confirm stability before proceeding.
>
> **RE-CHECK (2026-07-10):** The previous dispatch of this plan did NOT apply
> either pending phase — `tests/backup_test.go:28` still calls `src.Conn(nil)`
> (verified: `go test ./...` still panics with the exact deadlock described in
> Phase A, after a 90s timeout) and `tinywasm/sqlite/driver/` does not exist
> (verified: `tinywasm/sqlite/go.mod` still depends on `modernc.org/sqlite`
> directly, and `git status` in that repo is clean — no migration was ever
> run). `scripts/migrate_to_tinywasm.sh` and `driver/README.md` DO already
> exist in this repo, so Phase B's Step 1 and Step 3 are done — only the
> actual dry-run + real migration (Steps 2, 4-6) remain. Start again at
> Phase A; do not skip it because the script/README already exist.

---

## Execution Roadmap

| Phase | File | Nature | Status |
|-------|------|--------|--------|
| **A** | [A_FIX_TESTS.md](A_FIX_TESTS.md) | Fix `TestBackup` deadlock | 🔲 Pending (Start here!) |
| **B** | [8_MIGRATION.md](8_MIGRATION.md) | Migration script dry-run + copy to `tinywasm/sqlite` | 🔲 Pending (script + README already exist; run Steps 2, 4-6) |

---

## Context: What Was Already Done

| Phase | Result |
|-------|--------|
| 3, 4, 5 | Complete — archived in [PLAN_COMPLETED.md](PLAN_COMPLETED.md) |
| 6 | Skipped — `modernc.org/libc` accepted as permanent dep ([6_DEPS_LIBC.md](6_DEPS_LIBC.md)) |
| 7 | Deps clean ✅ — `go.mod` has one direct dep (`modernc.org/libc`), `go build ./...` passes ✅ |
| 8 (partial) | `scripts/migrate_to_tinywasm.sh` and `driver/README.md` exist, but were never run against the real `tinywasm/sqlite` repo |

---

## Phase A: Fix `TestBackup` Deadlock

- **File:** [A_FIX_TESTS.md](A_FIX_TESTS.md)
- **Problem:** `go test ./...` panics with a deadlock in `tests/backup_test.go:28`.
  `db.Conn()` is called after `db.Close()`, causing `database/sql` to panic.
- **Status:** 🔲 Pending — **NOT applied**. `tests/backup_test.go:28` still has
  `src.Conn(nil)` with no `"context"` import. Confirmed reproducing the exact
  panic trace from `A_FIX_TESTS.md` on 2026-07-10.

## Phase B: Migration Script & Final Validation

- **File:** [8_MIGRATION.md](8_MIGRATION.md)
- **Goal:** Run `scripts/migrate_to_tinywasm.sh` dry-run, then copy `driver/` to
  `tinywasm/sqlite/driver/` and validate with `go build ./...` + `go test ./...`.
- **Status:** 🔲 Pending — Step 1 (write script) and Step 3 (write README) are
  already done. Steps 2 (dry-run), 4 (final local `go test ./...`, blocked on
  Phase A), 5 (run migration into `/home/cesar/Dev/Project/tinywasm/sqlite`),
  and 6 (verify in target repo) were never executed —
  `/home/cesar/Dev/Project/tinywasm/sqlite/driver/` does not exist and that
  repo's `go.mod` still depends directly on `modernc.org/sqlite`, not the
  migrated driver package.

---

## Development Rules

- **Max 500 lines per file** (new hand-written files only).
- **No external assertion libraries.** Standard `testing` package only.
- **No global state.** Dependency injection via interfaces.
- Coverage target: **≥ 90%** (verify with `go test -cover`).
