# PLAN: Driver Encapsulation & Dependency Cleanup

## Purpose

The end goal is to finalize the self-contained `driver/` package so it can be copied verbatim into `tinywasm/sqlite/driver/`. The `go.mod` must end up 100% clean, meaning it has zero `modernc.org/*` or other third-party dependencies outside of the Go standard library and `tinywasm/*` packages.

> **IMPORTANT FOR THE AGENT:** Execute all pending phases in sequential order. For each phase, read its corresponding markdown document in the `docs/` directory for detailed instructions. After completing a phase, run `go test ./...` to ensure stability and proceed to the next phase. Completion is reached once Phase 8 is successfully finished.
> Note: Completed phases (Phases 3, 4, and 5) have been archived to [PLAN_COMPLETED.md](PLAN_COMPLETED.md).

## Execution Roadmap for Pending Phases

| Phase | Script / File | Nature | Target | Status |
|-------|---------------|--------|--------|--------|
| **6** | [6_DEPS_LIBC.md](6_DEPS_LIBC.md) | Inline `modernc.org/libc` into `driver/libc/` | ✅ Clean libc | 🔲 Pending |
| **7** | [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md) | Remove remaining `modernc.org` and other deps | ✅ Standard lib only | 🔲 Pending |
| **8** | [8_MIGRATION.md](8_MIGRATION.md) | Migration script + dry-run | ✅ Deployment ready | 🔲 Pending |

---

## Pending Phases Details

### Phase 6: Inline `modernc.org/libc`
- **Goal:** Eliminate `modernc.org/libc` from `go.mod` by inlining its source flat into `driver/libc/`.
- **Reference:** [6_DEPS_LIBC.md](6_DEPS_LIBC.md)
- **Status:** 🔲 Pending (Start here!)

### Phase 7: Clean Remaining External Dependencies
- **Goal:** Remove `github.com/google/uuid`, `github.com/dustin/go-humanize`, `github.com/mattn/go-isatty`, `github.com/ncruces/go-strftime`, `github.com/remyoudompheng/bigfft`, and `golang.org/x/exp` via inlining or replacement.
- **Reference:** [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md)
- **Status:** 🔲 Pending

### Phase 8: Migration Script & Final Validation
- **Goal:** Write & validate `scripts/migrate_to_tinywasm.sh`, perform dry-runs, and prepare final copy to `tinywasm/sqlite`.
- **Reference:** [8_MIGRATION.md](8_MIGRATION.md)
- **Status:** 🔲 Pending

---

## Development Rules

- **Max 500 lines per file** (applies only to new hand-written files, not auto-generated sources).
- **No external assertion libraries.** Standard `testing` package only.
- **No global state.** Use dependency injection via interfaces.
- Coverage target: **≥ 90%** (verify with `go test -cover`).
