# Phase 1A: Move Root `.go` Files into `driver/`

> **Master Plan:** [PLAN.md](PLAN.md)
> **Phase index:** [1_DRIVER_ORGANIZE.md](1_DRIVER_ORGANIZE.md)
> **Previous:** [0_MODULE_CLEANUP.md](0_MODULE_CLEANUP.md) ← must be ✅ complete
> **Next:** [1B_MOVE_LIB.md](1B_MOVE_LIB.md)

## Prerequisites

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

---

## Goal

Move all `.go` files from the repo root into `driver/`. Do **not** change any
file content in this sub-stage. Do **not** touch `lib/`, `vfs/`, or `vtab/` yet.

---

## Steps

### Step 1 — Create `driver/` directory

```bash
mkdir -p driver
```

### Step 2 — Move root `.go` files

```bash
git mv *.go driver/
```

`go.mod`, `go.sum`, and `README.md` stay at the root — git will skip them
automatically.

### Step 3 — Rename `mutex.go` to avoid stdlib clash

```bash
git mv driver/mutex.go driver/sqlite_mutex.go
```

### Step 4 — Commit

```bash
git add -A
git commit -m "move(1a): root *.go → driver/"
```

---

## Acceptance Criteria

| Criterion | Check |
|-----------|-------|
| No `.go` files remain at repo root (only `go.mod`, `go.sum`, `README.md`) | 🔲 |
| `driver/sqlite_mutex.go` exists (no `driver/mutex.go`) | 🔲 |
| `lib/`, `vfs/`, `vtab/` remain at root (unchanged) | 🔲 |
| Commit created | 🔲 |
