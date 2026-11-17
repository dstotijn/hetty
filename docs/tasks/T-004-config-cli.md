---
id: T-004
title: Config, CLI flags/env/TOML, app wiring, structured logging
milestone: M0
role: backend
status: todo
depends_on: [T-002]
owner: ""
branch: task/T-004-config-cli
paths:
  - cmd/ratatoskr/**
  - internal/config/**
  - internal/app/**
  - internal/logging/**
  - go.mod
  - go.sum
spec:
  - docs/01-feature-spec.md#i-cli-f-cli
  - docs/03-architecture.md §2, §9
  - docs/05-go-conventions.md §1, §3
---

## Goal
`ratatoskr` starts from a validated `config.Config` built from flags, `RATATOSKR_*` env vars and `<data-dir>/config.toml` (precedence flags > env > file > defaults),
wires an `app.App` with graceful shutdown, and logs via slog to stderr + an in-memory ring buffer.

## Acceptance criteria
- [ ] Flags from F-CLI exist (unimplemented ones return a clear "not implemented yet" error), `ratatoskr version` prints version/commit.
- [ ] Precedence covered by table tests.
- [ ] SIGINT/SIGTERM triggers graceful shutdown within 5 s (test with `signal.NotifyContext`, synctest where possible).
- [ ] Ring buffer keeps last 10k lines; `slog.NewMultiHandler` used.
- [ ] Data dir created with 0700 and accessed via `os.Root`.
- [ ] `--debug` mounts `/debug/pprof` (incl. goroutineleak) on the UI server only.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
