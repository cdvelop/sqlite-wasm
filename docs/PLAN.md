# PLAN: Driver Encapsulation & Dependency Cleanup

## Purpose

The end goal is to finalize the self-contained `driver/` package so it can be copied verbatim into `tinywasm/sqlite/driver/`. The `go.mod` must end up 100% clean, meaning it has zero `modernc.org/*` or other third-party dependencies outside of the Go standard library and `tinywasm/*` packages.

> **IMPORTANT FOR JULES:** Execute all phases listed below in order. For each phase, read its corresponding markdown document for detailed instructions. After completing a phase, run `go test ./...` to ensure stability and then proceed immediately to the next phase. Completion is reached once Phase 8 is successfully finished.

## Execution Roadmap

| Phase | Script / File | Nature | Target |
|-------|---------------|--------|--------|
| **3** | [3_TESTS_MOVE.md](3_TESTS_MOVE.md) | Move tests to `tests/` (DONE) | ✅ Tests in `tests/` |
| **4** | [4_TESTS_DOMAIN.md](4_TESTS_DOMAIN.md) | Domain test split; coverage ≥ 90% | ✅ Coverage ≥ 90% |
| **5** | [5_DEPS_SMALL.md](5_DEPS_SMALL.md) | Inline `modernc.org/mathutil` and `modernc.org/memory` | ✅ Localized deps |
| **6** | [6_DEPS_LIBC.md](6_DEPS_LIBC.md) | Inline `modernc.org/libc` into `driver/libc/` | ✅ Clean libc |
| **7** | [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md) | Remove remaining `modernc.org` and other deps | ✅ standard lib only |
| **8** | [8_MIGRATION.md](8_MIGRATION.md) | Migration script + dry-run | ✅ Deployment ready |

---

## Pending Phases

| File | Phase | Goal | Status |
|------|-------|------|--------|
| [3_TESTS_MOVE.md](3_TESTS_MOVE.md) | 3 | Move all tests to `tests/`; add build tags; `go test ./...` passes. **Partially done:** `tests/setup_test.go` and `tests/vfs_test.go` already exist. Complete remaining test file moves. | ✅ Complete |
| [4_TESTS_DOMAIN.md](4_TESTS_DOMAIN.md) | 4 | Subdivide tests by domain (conn, stmt, vfs, vtab, backup); coverage ≥ 90% | ✅ Complete |
| [5_DEPS_SMALL.md](5_DEPS_SMALL.md) | 5 | Inline `modernc.org/mathutil` and `modernc.org/memory` into `driver/mathutil/` and `driver/memory/`. `modernc.org/fileutil` is NOT in `go.mod` — skip it. | ✅ Complete |
| [6_DEPS_LIBC.md](6_DEPS_LIBC.md) | 6 | Inline `modernc.org/libc` (including `libc/sys/types`) into `driver/libc/` — one subdirectory level only, no deeper nesting. | 🔲 Pending |
| [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md) | 7 | Remove `github.com/google/uuid`, `github.com/dustin/go-humanize`, `github.com/mattn/go-isatty`, `github.com/ncruces/go-strftime`, `github.com/remyoudompheng/bigfft`, `golang.org/x/exp` via inline or `go mod tidy`. Evaluate `golang.org/x/sys` — keep only if unavoidable. | 🔲 Pending |
| [8_MIGRATION.md](8_MIGRATION.md) | 8 | Write & validate `scripts/migrate_to_tinywasm.sh`; dry-run copy into `tinywasm/sqlite`. | 🔲 Pending |

---

## Development Rules

- **Max 500 lines per file** (applies only to new hand-written files, not auto-generated sources).
- **No external assertion libraries.** Standard `testing` package only.
- **No global state.** Use dependency injection via interfaces.
- Coverage target: **≥ 90%** from Phase 4 onward.
