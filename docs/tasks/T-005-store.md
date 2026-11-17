---
id: T-005
title: SQLite store - open, WAL, migrations, write serializer, schema v1
milestone: M0
role: backend
status: todo
depends_on: [T-002]
owner: ""
branch: task/T-005-store
paths:
  - internal/store/**
  - go.mod
  - go.sum
spec:
  - docs/03-architecture.md §5
  - docs/adr/0002-sqlite-storage.md
---

## Goal
`internal/store` opens instance and project databases with the required pragmas, runs embedded forward-only migrations, provides a
serialized writer + read pool, transaction helpers, and the v1 schemas (tables from docs/03 §5, created even if unused yet).

## Acceptance criteria
- [ ] `store.OpenInstance(ctx, root *os.Root)` / `store.OpenProject(ctx, root)` return handles; pragmas verified by test.
- [ ] Migration runner is idempotent, records version, refuses to open a DB with a newer schema (clear error).
- [ ] Concurrent writers (100 goroutines × 1k inserts) complete without `SQLITE_BUSY` errors (`-race`).
- [ ] Batch insert helper for requests achieves ≥ 20k rows/s on CI hardware (benchmark).
- [ ] SQLite user functions registered: `regexp(pattern, text)`, `bcontains_fold(blob, needle)` with tests.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
