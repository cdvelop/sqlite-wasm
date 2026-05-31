# PLAN: Driver Encapsulation & Dependency Cleanup

## Purpose

The end goal is to finalize the self-contained `driver/` package so it can be copied verbatim into `tinywasm/sqlite/driver/`. The `go.mod` must keep only `modernc.org/libc` and `tinywasm/*` packages as direct dependencies. All other third-party dependencies must be eliminated.

> **Decision (2026-05-31):** Inlining `modernc.org/libc` was attempted in Phase 6 but deferred — the library has 2000+ platform-specific files and is too large to inline safely. It is accepted as a permanent direct dependency. See [6_DEPS_LIBC.md](6_DEPS_LIBC.md) for the documented decision.

> **IMPORTANT FOR THE AGENT:** Execute all pending phases in sequential order. For each phase, read its corresponding markdown document in the `docs/` directory for detailed instructions. After completing a phase, run `go test ./...` to ensure stability and proceed to the next phase. Completion is reached once Phase 8 is successfully finished.
> Note: Completed phases (Phases 3, 4, and 5) have been archived to [PLAN_COMPLETED.md](PLAN_COMPLETED.md). Phase 6 has been formally skipped — start at Phase 7.

## Execution Roadmap for Pending Phases

| Phase | Script / File | Nature | Target | Status |
|-------|---------------|--------|--------|--------|
| **6** | [6_DEPS_LIBC.md](6_DEPS_LIBC.md) | ~~Inline `modernc.org/libc`~~ — **SKIPPED** (accepted dep) | `modernc.org/libc` stays | ✅ Decided |
| **7** | [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md) | Remove all deps NOT required by `modernc.org/libc` | Only `modernc.org/libc` + its transitives | 🔲 Pending (Start here!) |
| **8** | [8_MIGRATION.md](8_MIGRATION.md) | Migration script + dry-run | ✅ Deployment ready | 🔲 Pending |

---

## Pending Phases Details

### Phase 6: ~~Inline `modernc.org/libc`~~ — SKIPPED
- **Decision:** `modernc.org/libc` is accepted as a permanent direct dependency. 2000+ platform-specific files make inlining impractical and a maintenance burden.
- **Reference:** [6_DEPS_LIBC.md](6_DEPS_LIBC.md) (documents the decision)
- **Status:** ✅ Decided — skip, proceed to Phase 7

### Phase 7: Clean Remaining External Dependencies
- **Goal:** Remove all deps that are NOT `modernc.org/libc` or its transitive requirements. After this phase `go.mod` contains only `modernc.org/libc` (direct) plus whatever it pulls in indirectly via `go mod tidy`; no other hand-authored direct deps remain.
- **Reference:** [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md)
- **Status:** 🔲 Pending (Start here!)

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
