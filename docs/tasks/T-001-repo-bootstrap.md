---
id: T-001
title: Bootstrap the fork - module path, Go 1.27, internal layout, Makefile, lint, CI
milestone: M0
role: lead
status: in_review
depends_on: []
owner: "claude-code (claude-opus-5-5)"
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
  # Widened by lead (T-001): the module rename and "keep Hetty code compiling"
  # require touching legacy Go code and the gqlgen config.
  - gqlgen.yml
  - tools.go
  - pkg/**
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
- [x] `go build ./...` and `go vet ./...` pass with Go 1.27.
- [ ] `make check` passes locally and in CI on ubuntu/macos/windows. (local: green on darwin/arm64; CI not run yet, branch not pushed)
- [x] `go fix -diff ./...` is empty.
- [x] README credits Hetty and links the docs.

## Notes from lead
Replace `OWNER` with the real GitHub owner (ask the operator). Keep Hetty code compiling in `pkg/` for now.

## Plan (implementer, ≤15 lines)
1. go.mod: module github.com/pocikode/ratatoskr, go 1.27, toolchain go1.27.4; rewrite Hetty import paths; go mod tidy.
2. Replace tools.go with `tool` directives (gqlgen, govulncheck); keep pkg/api generated code compiling.
3. Apply `go fix ./...` so `go fix -diff` is empty (modernizer rewrites in legacy pkg/).
4. Makefile: gen, check, test-int, build, fuzz, run (+ web-* no-ops until T-003).
5. .golangci.yml v2; legacy pkg/, cmd/hetty excluded with T-002 comment if needed.
6. CI: replace Hetty workflows with ci.yml matrix ubuntu/macos/windows running `make check`.
7. LICENSE second copyright line; short README; cmd/ratatoskr/main.go (version); internal/doc.go.
8. go build/vet/fix -diff + make check green; Handoff.

## Out-of-scope change needed (if any)
- Resolved by lead (this is a lead task): the module rename rewrites imports in `pkg/**` and `cmd/hetty/**`,
  `go fix` and the gqlgen bump touch `pkg/**`, `tools.go` is replaced by `tool` directives, `gqlgen.yml` holds the
  module path. `paths` widened above with a note (`cmd/**` was already partially covered via `cmd/ratatoskr/**`;
  `cmd/hetty/**` edits are limited to the import rewrite and a build tag).
- Not done, left to T-002 (paths added there): `.goreleaser.yml`, `Dockerfile`, `.dockerignore` still build
  `cmd/hetty` + `admin/` and reference `dstotijn/hetty`.

## Handoff (implementer)
**What changed**
- `go.mod`: `module github.com/pocikode/ratatoskr`, `go 1.27`, `toolchain go1.27.4`; `tool` directives for
  `github.com/99designs/gqlgen` (bumped v0.14.0 → v0.17.98) and `golang.org/x/vuln/cmd/govulncheck`; `tools.go` removed.
- Import paths rewritten in `pkg/**`, `cmd/hetty/**`, `gqlgen.yml`. `pkg/api/generated.go` + `models_gen.go`
  regenerated with `go tool gqlgen generate` (hand-written `resolvers.go` kept; gqlgen's rewrite of it discarded).
- `cmd/hetty/*.go` gated behind `//go:build legacy_hetty`: it `//go:embed`s the Next.js admin build that no
  longer exists, so `go build ./...` failed even upstream. T-002 deletes it.
- `go fix ./...` applied to legacy code (ioutil, loop var copies, `wg.Go`, `any`).
- Legacy bug fixes needed for `go test -race`: shared `math/rand.Rand` ULID entropy (not goroutine-safe) replaced
  by `crypto/rand.Reader` in `pkg/proj`, `pkg/proxy`, `pkg/sender` and four test files.
- Legacy flaky tests (fixed `time.Sleep(10ms)` waits; 3–10 % failure rate upstream, with or without `-race`):
  `pkg/proxy/intercept` tests now poll `ItemByID`; the two "…that's done" subtests asserted `ErrRequestDone`, which
  only exists in a scheduler-dependent window between cancel and queue removal. They now wait for the intercept
  goroutine to finish and assert the settled state `ErrRequestNotFound` (lead decision, see commit 69c59d3).
  `pkg/reqlog` test polls until the response log is stored. 0 failures in 200 `-race` iterations.
- New: `cmd/ratatoskr/main.go` (prints `ratatoskr <version>`, `-X main.version`), `internal/doc.go`, `Makefile`
  (gen, check = fix-check+vet+lint+test+web-check+web-test, test-int, fuzz, vulncheck, build, run, clean; `web-*`
  run `bun` only when `web/package.json` exists), `.golangci.yml` (v2; linters from docs/05 §3 + depguard
  denying `net/http/httputil` per ADR-0003 + nolintlint requiring explanations), `.editorconfig`, `.gitignore` entries.
- CI: Hetty `build-test.yml`/`lint.yml` and `FUNDING.yml` removed; `.github/workflows/ci.yml` runs `make check` on
  ubuntu/macos/windows (bash shell, LF checkout, golangci-lint v2.15.0 via `go install`, `make` via choco on Windows).
  Issue template discussion link points to our repo.
- `LICENSE`: `Copyright (c) 2026 pocikode` added below Hetty's line. `README.md` rewritten.

**How to verify**
`go build ./... && go vet ./... && go fix -diff ./... && make check && make build && ./ratatoskr`
(expect `ratatoskr <git describe>`). `go tool govulncheck ./...` → no vulnerabilities. `GOOS=windows go vet ./...` ok.

**Limitations**
- CI matrix not yet exercised (nothing pushed). Windows assumptions: `make` installable via choco, gcc on PATH for
  `-race` (cgo). If either fails, fix in `.github/workflows/ci.yml`.
- Legacy `pkg/` and `cmd/hetty/` are excluded from golangci-lint (89 issues, mostly revive) with a T-002 comment.
- `make gen` runs `go generate ./...` (no directives yet; T-007 adds gqlgen's). `make run` passes `--data-dir`,
  which `cmd/ratatoskr` ignores until T-004.

**Follow-ups**
- T-002: delete `cmd/hetty` (drop `legacy_hetty` tag), drop the `pkg/`/`cmd/hetty/` lint exclusions as code moves,
  fix or delete `.goreleaser.yml`/`Dockerfile`/`.dockerignore` (paths added to T-002).
- Consider a `govulncheck` CI job once legacy deps are gone.

## Review (reviewer)
