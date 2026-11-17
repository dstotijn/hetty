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
The forked repo builds with Go 1.27 under the new module path, has the target directory skeleton, and CI runs `make check` on Linux and macOS.
(Operator decision 2026-11-17: Windows is not a supported target; Unix-like only.)

## Scope
- In: `go.mod` → `module github.com/OWNER/ratatoskr`, `go 1.27`, `toolchain go1.27.x`; update imports; `tool` directives (gqlgen, govulncheck);
  Makefile targets from docs/05 §5 (web targets may be no-ops until T-003); `.golangci.yml` (v2 format);
  GitHub Actions matrix; LICENSE: keep Hetty line, add `Copyright (c) 2026 <owner>`; README rewritten (short, our name, credits Hetty);
  empty `cmd/ratatoskr/main.go` that prints version; `internal/doc.go`.
- Out: deleting Hetty code (T-002), any feature work.

## Acceptance criteria
- [x] `go build ./...` and `go vet ./...` pass with Go 1.27.
- [ ] `make check` passes locally and in CI on ubuntu/macos. (local: green on darwin/arm64; CI not run yet, branch not pushed)
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
6. CI: replace Hetty workflows with ci.yml matrix ubuntu/macos running `make check`.
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
  ubuntu/macos (golangci-lint v2.15.0 via `go install`). Windows dropped by operator decision; Makefile has no `.exe` handling.
  Issue template discussion link points to our repo.
- `LICENSE`: `Copyright (c) 2026 pocikode` added below Hetty's line. `README.md` rewritten.

**How to verify**
`go build ./... && go vet ./... && go fix -diff ./... && make check && make build && ./ratatoskr`
(expect `ratatoskr <git describe>`). `go tool govulncheck ./...` → no vulnerabilities.

**Limitations**
- CI matrix not yet exercised (nothing pushed).
- Legacy `pkg/` and `cmd/hetty/` are excluded from golangci-lint (89 issues, mostly revive) with a T-002 comment.
- `make gen` runs `go generate ./...` (no directives yet; T-007 adds gqlgen's). `make run` passes `--data-dir`,
  which `cmd/ratatoskr` ignores until T-004.

**Follow-ups**
- T-002: delete `cmd/hetty` (drop `legacy_hetty` tag), drop the `pkg/`/`cmd/hetty/` lint exclusions as code moves,
  fix or delete `.goreleaser.yml`/`Dockerfile`/`.dockerignore` (paths added to T-002).
- Consider a `govulncheck` CI job once legacy deps are gone.

## Review (reviewer)
Reviewed by claude-code (claude-opus-5-5), 2026-11-17, against `7c05016...1bd5efe` (merged as PR #1, `1676833`).
`make check` on merged main: green (exit 0). `go build ./...` ok. CI green on ubuntu/macos (operator).
AC: build/vet ✔, `go fix -diff` empty ✔, `make check` local + CI ✔, README credits Hetty and links existing docs ✔.
Hard rules: no httputil added (depguard now denies it), no generated code hand-edited (gqlgen regenerated), no secrets.

- [ ] [minor] pkg/proxy/intercept/intercept_test.go:82,203 — the `ErrRequestDone` branch of `ModifyRequest`/`ModifyResponse`
  is no longer exercised by any test (assertions moved to the settled `ErrRequestNotFound` state). Acceptable for legacy
  code, but the behaviour must not be lost → when `internal/intercept` is written, cover "modify after cancel" deterministically
  with `testing/synctest` (docs/05) instead of sleeps; add this to the intercept task's acceptance criteria.
- [ ] [minor] docs/tasks/T-001-repo-bootstrap.md:10 — `cmd/hetty/**` (4 files modified: import rewrite + `legacy_hetty` build tag)
  is not in `paths`; only `cmd/ratatoskr/**` is. The Out-of-scope note explains it, but `paths` should match the diff →
  add `cmd/hetty/**` to the widened list.
- [ ] [minor] .github/workflows/ci.yml — CI runs only `make check`; docs/05 §2 requires `govulncheck`, and `make build`
  (ldflags, `-trimpath`, CGO_ENABLED=0) is never exercised in CI → add `make vulncheck` (once legacy deps are gone, per
  Handoff follow-up) and `make build` steps, or track them in a task.
- [ ] [nit] pkg/proj/proj.go:23, pkg/proxy/proxy.go:25, pkg/sender/sender.go:22 (+4 test files) — `//nolint:gosec` left above
  `var ulidEntropy = rand.Reader`; it was for `math/rand` and is now stale and lacks the explanation nolintlint requires
  (masked only because `pkg/` is lint-excluded) → drop the directives (T-002 can do it).
- [ ] [nit] Goal says "has the target directory skeleton"; only `internal/doc.go` exists. Fine given Scope, as later tasks create
  their packages; consider rewording the Goal so it doesn't read as unmet.

**Verdict: approve.** No blocker/major findings. Already merged; operator may set `status: done`. Minor items can go to
T-002 or a follow-up task.
