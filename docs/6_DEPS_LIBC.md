# Phase 6: `modernc.org/libc` — Decision: Accepted Dependency

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [5_DEPS_SMALL.md](5_DEPS_SMALL.md) ← complete
> **Next:** [7_DEPS_CLEAN.md](7_DEPS_CLEAN.md)

## Status: SKIPPED — `modernc.org/libc` is an accepted permanent dependency

---

## Decision (2026-05-31)

Inlining `modernc.org/libc` into `driver/libc/` was attempted by an automated
agent and abandoned. The library exceeds 2000 source files with heavy
platform-specific build tags (`linux/amd64`, `darwin/arm64`, `windows/amd64`,
etc.). Maintaining an inlined copy would be an ongoing burden with every libc
upstream update.

The Decision Gate defined in the original plan explicitly covered this scenario:

> *"If > 500 files: consider deferring to a separate dedicated effort and document
> `modernc.org/libc` as an accepted remaining dependency."*

**`modernc.org/libc` remains as a direct dependency in `go.mod`.**

---

## Impact on go.mod

The following entries are expected to remain after this decision (managed
automatically by `go mod tidy`):

```
modernc.org/libc        — direct (SQLite CGo-free runtime)
modernc.org/mathutil    — indirect (required by libc)
modernc.org/memory      — indirect (required by libc)
```

Whether additional transitive deps from libc (e.g. `golang.org/x/sys`,
`github.com/remyoudompheng/bigfft`) remain is determined by Phase 7 auditing.

---

## No action required

Proceed directly to **[Phase 7](7_DEPS_CLEAN.md)**.
