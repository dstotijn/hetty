---
id: T-001
title: Bootstrap the fork - module path, Go 1.27, internal layout, Makefile, lint, CI
milestone: M0
role: lead
status: todo
depends_on: []
owner: ""
branch: task/T-001-repo-bootstrap
paths:
  - go.mod
  - go.sum
  - Makefile
  - .golangci.yml
  - .github/**
  - LICENSE
  - README.md
  - .gitignore
  - .editorconfig
  - internal/doc.go
  - cmd/ratatoskr/**
  - docs/**
spec:
  - docs/03-architecture.md §2, §3
  - docs/05-go-conventions.md §2, §5
  - docs/adr/0001-fork-hetty-strangler-rewrite.md
---

## Goal
The forked repo builds with Go 1.27 under the new module path, has the target directory skeleton, and CI runs `make check` on Linux, macOS and Windows.

## Scope
- In: `go.mod` → `module github.com/OWNER/ratatoskr`, `go 1.27`, `toolchain go1.27.x`; update imports; `tool` directives (gqlgen, govulncheck);
  Makefile targets from docs/05 §5 (web targets may be no-ops until T-003); `.golangci.yml` (v2 format);
  GitHub Actions matrix; LICENSE: keep Hetty line, add `Copyright (c) 2026 <owner>`; README rewritten (short, our name, credits Hetty);
  empty `cmd/ratatoskr/main.go` that prints version; `internal/doc.go`.
- Out: deleting Hetty code (T-002), any feature work.

## Acceptance criteria
- [ ] `go build ./...` and `go vet ./...` pass with Go 1.27.
- [ ] `make check` passes locally and in CI on ubuntu/macos/windows.
- [ ] `go fix -diff ./...` is empty.
- [ ] README credits Hetty and links the docs.

## Notes from lead
Replace `OWNER` with the real GitHub owner (ask the operator). Keep Hetty code compiling in `pkg/` for now.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
