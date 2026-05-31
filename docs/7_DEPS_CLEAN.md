# Phase 7: Clean Remaining External Dependencies

> **Master Plan:** [PLAN.md](PLAN.md)
> **Previous:** [6_DEPS_LIBC.md](6_DEPS_LIBC.md) ← Phase 6 skipped, `modernc.org/libc` accepted
> **Next:** [8_MIGRATION.md](8_MIGRATION.md)

## Goal

Remove all direct dependencies from `go.mod` that are **not** `modernc.org/libc`,
`tinywasm/*` packages, or transitive deps pulled in by libc. After this phase
`go.mod` must have only these permitted direct deps:

- `modernc.org/libc` (accepted CGo-free SQLite runtime)
- `tinywasm/*` packages (ecosystem deps — always permitted)

Indirect entries are managed by `go mod tidy` and are acceptable as long as they
originate from the permitted deps above.

---

## Current `go.mod` state (before this phase)

```
require (
    github.com/remyoudompheng/bigfft        v0.0.0-...
    golang.org/x/sys                        v0.37.0
    modernc.org/libc                        v1.67.6
)

require (
    github.com/dustin/go-humanize           v1.0.1   // indirect
    github.com/google/uuid                  v1.6.0   // indirect
    github.com/mattn/go-isatty              v0.0.20  // indirect
    github.com/ncruces/go-strftime          v1.0.0   // indirect
    golang.org/x/exp                        v0.0.0-... // indirect
    modernc.org/mathutil                    v1.7.1   // indirect
    modernc.org/memory                      v1.11.0  // indirect
)
```

---

## Step 1 — Audit which deps `driver/` actually imports directly

```bash
grep -rh '"' driver/*.go driver/vfs/*.go driver/vtab/*.go 2>/dev/null \
    | grep -oP '"[^"]*"' | sort -u \
    | grep -v 'modernc.org/libc\|^"[a-z]' \
    | grep -v 'github.com/cdvelop/sqlite-wasm'
```

For each dep found, check whether it is also imported by libc itself:

```bash
GOMOD=$(go env GOMODCACHE)
LIBC_VER=$(grep "modernc.org/libc" go.mod | awk '{print $2}')
grep -r "dustin\|go-isatty\|strftime\|bigfft\|google/uuid\|golang.org/x/exp" \
    "$GOMOD/modernc.org/libc@$LIBC_VER" --include="*.go" -l
```

Deps that appear in libc source are transitive — `go mod tidy` handles them.
Deps that appear **only** in `driver/` code are direct and must be removed.

---

## Step 2 — Remove each direct dep not needed by libc

For each dep confirmed as NOT used by `modernc.org/libc`:

### `github.com/google/uuid`

```bash
grep -r "google/uuid" driver/ --include="*.go"
```

If used: replace with `github.com/tinywasm/unixid` (permitted ecosystem dep):

```bash
go get github.com/tinywasm/unixid
```

```go
import "github.com/tinywasm/unixid"

idHandler, err := unixid.NewUnixID()
if err != nil {
    return fmt.Errorf("id init: %w", err)
}
id := idHandler.GetNewID()
```

| Old (`google/uuid`) | New (`tinywasm/unixid`) |
|---------------------|------------------------|
| `uuid.New().String()` | `idHandler.GetNewID()` |
| `uuid.UUID` field type | `string` |
| `uuid.Must(uuid.Parse(s))` | `idHandler.Validate(s)` |

### `github.com/dustin/go-humanize`

```bash
grep -r "go-humanize\|humanize\." driver/ --include="*.go"
```

If used (typically `humanize.Bytes(n)`): replace inline:

```go
func humanBytes(n uint64) string {
    switch {
    case n >= 1<<30: return fmt.Sprintf("%.1f GB", float64(n)/(1<<30))
    case n >= 1<<20: return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
    case n >= 1<<10: return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
    default:         return fmt.Sprintf("%d B", n)
    }
}
```

### `github.com/mattn/go-isatty`

```bash
grep -r "go-isatty\|isatty\." driver/ --include="*.go"
```

If used (typically `isatty.IsTerminal(fd)`): replace with stdlib `syscall`:

```go
import "syscall"

func isTerminal(fd uintptr) bool {
    _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, 0)
    return err == 0
}
```

### `github.com/ncruces/go-strftime`

```bash
grep -r "go-strftime\|strftime\." driver/ --include="*.go"
```

If used: replace with `time.Format` using the equivalent Go layout string.

### `github.com/remyoudompheng/bigfft`

```bash
grep -r "bigfft" driver/ --include="*.go"
```

If `driver/` code does NOT import it directly (it may be a transitive dep of
`modernc.org/mathutil` which libc uses), run `go mod tidy` after other removals —
it may disappear on its own or become an `// indirect` entry, which is acceptable.

### `golang.org/x/exp`

```bash
grep -r "golang.org/x/exp" driver/ --include="*.go"
```

If used: replace with stdlib equivalents (slices, maps packages in Go 1.21+).

### `golang.org/x/sys`

```bash
grep -r "golang.org/x/sys" driver/ --include="*.go"
```

If only in inlined/generated files or needed transitively by libc: keep it as
`// indirect` — this is acceptable. Do NOT remove if it breaks the build.

---

## Step 3 — Run `go mod tidy` and verify

```bash
go mod tidy
cat go.mod
```

Expected result — only one direct dep:

```
require modernc.org/libc vX.Y.Z

require (
    // indirect entries managed by go mod tidy — all must originate from modernc.org/libc
    ...
)
```

If any non-libc direct dep still appears, return to Step 2 and eliminate it.

---

## Step 4 — Build and test

```bash
go build ./...
go test ./...
```

Coverage must remain ≥ 90%.

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| `go.mod` direct deps are only `modernc.org/libc` and/or `tinywasm/*` | ✅ |
| No non-permitted direct deps remain (`google/*`, `dustin/*`, `mattn/*`, etc.) | ✅ |
| `go build ./...` succeeds | ✅ |
| `go test ./...` passes with ≥ 90% coverage | ✅ |
| `driver/README.md` documents `modernc.org/libc` as accepted dep | ✅ |
