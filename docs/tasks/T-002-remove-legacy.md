---
id: T-002
title: Remove Hetty admin UI and bbolt, keep the binary compiling with stubs
milestone: M0
role: backend
status: todo
depends_on: [T-001]
owner: ""
branch: task/T-002-remove-legacy
paths:
  - admin/**
  - pkg/**
  - cmd/**
  - internal/legacy/**
  - go.mod
  - go.sum
  - .golangci.yml
  - .goreleaser.yml
  - Dockerfile
  - .dockerignore
spec:
  - docs/03-architecture.md §3
---

## Goal
Legacy pieces we will never reuse are gone; the remaining Hetty packages that serve as references stay under `pkg/` untouched or are moved to `internal/legacy/` (read-only reference) so nothing in `internal/` imports them.

## Scope
- In: delete `admin/`, `pkg/db/bolt`, Next.js embed directives, `cmd/hetty`; move `pkg/proxy`, `pkg/filter`, `pkg/scope`, `pkg/proxy/intercept`, `pkg/chrome` to `internal/legacy/…` with a README stating "reference only, do not import".
  `cmd/ratatoskr` builds and starts an HTTP server on `127.0.0.1:8080` returning a placeholder page.
- Out: new features.

## Acceptance criteria
- [ ] No `admin/` dir, no bbolt dependency in `go.mod` (`go mod tidy`).
- [ ] `go build ./cmd/ratatoskr && ./ratatoskr` serves a placeholder on :8080.
- [ ] `internal/legacy/importguard_test.go` fails if any non-legacy package imports `internal/legacy` (uses `go list -deps -json`).
- [ ] `make check` green.

## Notes from lead (added in T-001)
- `cmd/hetty/*.go` carry `//go:build legacy_hetty` (the admin embed no longer builds); delete the package.
- `.golangci.yml` excludes `^pkg/` and `^cmd/hetty/` (comment references T-002). Narrow or remove the exclusions
  for whatever leaves `pkg/`; `internal/legacy/**` may be excluded the same way (reference-only code).
- `.goreleaser.yml`, `Dockerfile`, `.dockerignore` still build `cmd/hetty` + `admin/`: delete them or point them at
  `cmd/ratatoskr` (no release pipeline is specified yet; deleting is acceptable).

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
