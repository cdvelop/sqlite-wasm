# Phase 7: Clean Remaining External Dependencies

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [6_DEPS_LIBC.md](6_DEPS_LIBC.md) ← must be complete
> **Next:** [8_MIGRATION.md](8_MIGRATION.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Eliminate all remaining external dependencies that are not part of the `tinywasm/*`
ecosystem. After this phase, `go.mod` only contains:
- `github.com/tinywasm/unixid` (replaces `google/uuid`)
- Standard library (`golang.org/x/sys` is evaluated below)

---

## Remaining Deps After Phase 6

Run to get the current state:
```bash
cat go.mod
```

Likely candidates:

| Package | Action |
|---------|--------|
| `github.com/google/uuid` | ❌ Remove — replace with `github.com/tinywasm/unixid` |
| `github.com/dustin/go-humanize` | ❌ Remove — inline or eliminate usage |
| `github.com/mattn/go-isatty` | ❌ Remove — inline (it's tiny) or eliminate |
| `github.com/ncruces/go-strftime` | ❌ Remove — inline or use stdlib `time.Format` |
| `github.com/remyoudompheng/bigfft` | ❌ Remove — inline from `modernc.org/mathutil` context |
| `golang.org/x/exp` | ❌ Remove — replace usages with stdlib equivalents |
| `golang.org/x/sys` | ⚠️ Evaluate — may be required for low-level OS calls; keep if unavoidable |
| `github.com/google/pprof` | ❌ Remove — only a dev-profiling dep, not needed in the driver |

---

## Steps

### Step 1 — Audit actual usage (not just `go.mod`)

```bash
# Find which driver/ files actually USE each dep
grep -r "go-humanize\|go-isatty\|strftime\|bigfft\|google/uuid\|golang.org/x/exp" \
     driver/ --include="*.go"
```

If a dep appears only in `go.sum` but not used in `driver/*.go`, it was already
pulled transitively and will be removed automatically by `go mod tidy` after
previous phases. No action needed.

### Step 2 — Replace `google/uuid` → `tinywasm/unixid`

If any usage exists in `driver/*.go`:

```bash
grep -r "google/uuid\|uuid\.New\|uuid\.UUID" driver/ --include="*.go"
```

Replace:

| Old (`google/uuid`) | New (`tinywasm/unixid`) |
|---------------------|------------------------|
| `uuid.New().String()` | `idHandler.GetNewID()` |
| `uuid.UUID` field type | `string` |
| `uuid.Must(uuid.Parse(s))` | `idHandler.Validate(s)` |

Initialization:
```go
import "github.com/tinywasm/unixid"

idHandler, err := unixid.NewUnixID()
if err != nil {
    return fmt.Errorf("id init: %w", err)
}
id := idHandler.GetNewID()
```

Add to `go.mod`:
```bash
go get github.com/tinywasm/unixid
```

### Step 3 — Inline tiny utility packages

For packages with minimal usage (< 5 functions used):

**`go-isatty`** — usually just `isatty.IsTerminal(fd)`:
```go
// Replace with syscall-based check (linux/amd64 example):
func isTerminal(fd uintptr) bool {
    var t syscall.Termios
    _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd,
        syscall.TCGETS, uintptr(unsafe.Pointer(&t)), 0, 0, 0)
    return err == 0
}
```

**`go-humanize`** — usually `humanize.Bytes(n)`:
```go
// Replace with stdlib fmt:
func humanBytes(n uint64) string {
    switch {
    case n >= 1<<30: return fmt.Sprintf("%.1f GB", float64(n)/(1<<30))
    case n >= 1<<20: return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
    case n >= 1<<10: return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
    default: return fmt.Sprintf("%d B", n)
    }
}
```

### Step 4 — Evaluate `golang.org/x/sys`

```bash
grep -r "golang.org/x/sys" driver/ --include="*.go"
```

- If used only in the inlined libc files: accept it as a remaining dep (low-level OS interface).
- If used directly by our logic: replace with `syscall` stdlib equivalent.

Document the decision in `driver/README.md`:
> `golang.org/x/sys` is retained as an indirect dependency required by the
> embedded libc layer for OS-specific system call interfaces.

### Step 5 — Run `go mod tidy`

```bash
go mod tidy
cat go.mod
```

Verify only `tinywasm/*` (and optionally `golang.org/x/sys`) remain.

### Step 6 — Build and test

```bash
go build ./...
gotest
```

Coverage must remain ≥ 90%.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `google/uuid` absent from `go.mod` | ✅ |
| `go-humanize`, `go-isatty`, `ncruces/go-strftime`, `bigfft` absent | ✅ |
| `github.com/tinywasm/unixid` present if uuid was used | ✅ |
| `google/pprof`, `golang.org/x/exp` absent | ✅ |
| `golang.org/x/sys` decision documented | ✅ |
| `go build ./...` succeeds | ✅ |
| `gotest` passes with ≥ 90% coverage | ✅ |
