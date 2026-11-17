---
id: T-002
title: Remove Hetty admin UI and bbolt, keep the binary compiling with stubs
milestone: M0
role: backend
status: in_review
depends_on: [T-001]
owner: opencode/deepseek-v4.1-flash-1017
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
- [x] No `admin/` dir, no bbolt dependency in `go.mod` (`go mod tidy`).
- [x] `go build ./cmd/ratatoskr && ./ratatoskr` serves a placeholder on :8080.
- [x] `internal/legacy/importguard_test.go` fails if any non-legacy package imports `internal/legacy` (uses `go list -deps -json`).
- [x] `make check` green.

## Notes from lead (added in T-001)
- `cmd/hetty/*.go` carry `//go:build legacy_hetty` (the admin embed no longer builds); delete the package.
- `.golangci.yml` excludes `^pkg/` and `^cmd/hetty/` (comment references T-002). Narrow or remove the exclusions
  for whatever leaves `pkg/`; `internal/legacy/**` may be excluded the same way (reference-only code).
- `.goreleaser.yml`, `Dockerfile`, `.dockerignore` still build `cmd/hetty` + `admin/`: delete them or point them at
  `cmd/ratatoskr` (no release pipeline is specified yet; deleting is acceptable).

## Plan (implementer, ≤15 lines)
1. Delete `admin/`, `cmd/hetty/` (the `legacy_hetty` Next.js embed) and `pkg/db/bolt` (bbolt, ADR-0002).
2. `git mv pkg internal/legacy`: the whole remaining tree moves, not just the five named packages — `api`, `proj`,
   `reqlog`, `sender`, `log` import `proxy`/`filter`/`scope`/`proxy/intercept`, and the guard forbids `pkg/` →
   `internal/legacy/` imports. Reference tree stays intact; imports rewritten mechanically.
3. `internal/legacy/{doc.go,README.md}`: "reference only, do not import".
4. `internal/legacy/importguard_test.go`: `go list -deps -json -test ./...` from the module root; fail when a
   main-module package outside `internal/legacy/` imports it (test imports included).
5. `cmd/ratatoskr`: `-listen` (default `127.0.0.1:8080`), `-data-dir` (accepted, wired in T-004), `-version`;
   `http.Server` serving `GET /{$}` placeholder, ReadHeaderTimeout/IdleTimeout/MaxHeaderBytes, graceful shutdown.
6. Delete `.goreleaser.yml`, `Dockerfile`, `.dockerignore` (Hetty release/container infra for the deleted `cmd/hetty`).
7. `.golangci.yml`: drop `^pkg/` + `^cmd/hetty/` exclusions, exclude `^internal/legacy/` instead.
8. `go mod tidy` drops bbolt; verify `go build ./cmd/ratatoskr && ./ratatoskr` + `make check`.

## Out-of-scope change needed (if any)
Three stale references sit outside this task's `paths`:
- `gqlgen.yml` (root) still points at `pkg/api/schema.graphql`, `pkg/api/generated.go`, `pkg/api/resolvers.go` and the
  `github.com/pocikode/ratatoskr/pkg/api.ULID|URL` models. It is inert today (`go generate ./...` has no directives) and T-007
  replaces it with the `internal/api/graphql` config — delete it there rather than retargeting dead code now.
- `.gitignore` still lists `/cmd/hetty/admin`; the directory is gone.
- `AGENTS.md:62` says the knowledge graph "still contains legacy Hetty code (`admin/`, `pkg/`) until T-002 removes it" —
  refresh `graphify-out/` (`graphify update .`) and that line once this merges.

## Handoff (implementer)
**What changed**
- Deleted `admin/` (Next.js UI, replaced by `web/` in T-003), `cmd/hetty` (the `legacy_hetty` package with the admin `go:embed`
  directives) and `pkg/db/bolt` (bbolt → SQLite, ADR-0002). The two tests that only existed to drive the bolt repository
  (`pkg/reqlog/reqlog_test.go`, `pkg/sender/sender_test.go`) went with it; their `search_test.go` siblings stay.
- `pkg/{proxy,proxy/intercept,filter,scope,chrome,api,proj,reqlog,sender,log}` → `internal/legacy/…` via `git mv`, imports
  rewritten mechanically, no logic change. Only the five "reference" packages were named in Scope, but `api`, `proj`, `reqlog`,
  `sender` and `log` import them and the new guard forbids `pkg/` → `internal/legacy/` imports, so the reference tree moved as a
  unit rather than being deleted.
- New `internal/legacy/doc.go` + `README.md` ("reference only, do not import", with the §3 fate table) and
  `internal/legacy/importguard_test.go`: runs `go list -deps -json -test ./...` from the module root (found by walking up to
  `go.mod`), decodes with `encoding/json/v2` + `jsontext`, and fails on any main-module package outside `internal/legacy/` that
  imports it — test-only imports included.
- `cmd/ratatoskr`: `run`/`serve` split, `--listen` (default `127.0.0.1:8080`), `--data-dir` (accepted, used by T-004),
  `-version`; `GET /{$}` placeholder page, `ReadHeaderTimeout` 10 s, `IdleTimeout` 60 s, 1 MiB header cap, SIGINT/SIGTERM →
  `Shutdown` with a 5 s grace period, address logged after `net.Listen` so `--listen 127.0.0.1:0` is usable in parallel worktrees.
  Tests cover the handler (200/404/405/HEAD), the `serve` lifecycle on a real ephemeral listener, and flag parsing.
- `.golangci.yml`: `^pkg/` + `^cmd/hetty/` exclusions replaced by `^internal/legacy/` (linters and formatters). Stale
  `//nolint:gosec` above `ulidEntropy` removed in `proj`, `proxy`, `sender` and `proxy/intercept` (T-001 review nit).
- Deleted `.goreleaser.yml`, `Dockerfile`, `.dockerignore` (Hetty tap/scoop/snapcraft/ghcr config for the deleted `cmd/hetty`).
  `go mod tidy` dropped bbolt, go-cmp, ff/v3, go-homedir and truststore.

**How to verify**
- `make check` → `go fix -diff` empty, vet clean, `0 issues.` from golangci-lint, `go test -race ./...` green (web targets skip).
- `make build && ./ratatoskr` → `curl -i http://127.0.0.1:8080/` gives 200 + `text/html` placeholder, `/nope` 404, `POST /` 405,
  `Ctrl-C` logs `ratatoskr stopped`; `./ratatoskr -version` prints `ratatoskr <git describe>`.
- `grep -rn bbolt go.mod go.sum` → empty; `ls admin pkg` → neither exists.
- Guard: temporarily import `internal/legacy/filter` from any non-legacy package → `go test ./internal/legacy` fails with
  `<pkg> imports …/internal/legacy/filter` (verified, then reverted).
- `go tool govulncheck ./...` → no vulnerabilities (not part of `make check`).

**Known limitations**
- `--data-dir` is parsed and logged but not used; `internal/config` (T-004) owns it. The flag exists so `make run` keeps working.
- `internal/legacy` is still in `./...`: vet and `go test -race` cover it, golangci-lint does not. It is dead code until each
  rewrite task supersedes it.
- The placeholder has no auth/CSRF middleware, no logging middleware and no UI — those land with T-003/T-004/T-007.
- `internal/legacy/api/generated.go` is gqlgen output kept as-is (generated code, never edited by hand).

**Follow-ups (proposed tasks)**
- Delete `gqlgen.yml` with T-007.
- Delete each `internal/legacy/<pkg>` as its rewrite lands (T-010 rawhttp/proxy, T-014 query, T-018 scope, T-031 intercept,
  T-033 replay, T-041 files…), so the tree reaches zero; the guard keeps it from leaking back.
- Re-add a container/release pipeline when packaging is specified (no task yet).

## Review (reviewer)
Reviewer: claude-code/claude-opus-5-5 · round 1 · 2026-11-17

**Checks run**
- `make check`: green. `go fix -diff` empty, vet clean, golangci-lint `0 issues.`, `go test -race ./...` ok (`cmd/ratatoskr`,
  `internal/legacy` + its surviving `filter`, `proxy/intercept`, `reqlog`, `sender` tests); web targets skipped (no `web/` until T-003).
- `make test-int` not run: no proxy/store/rawhttp product code changed. `internal/legacy/proxy` is a pure move.
- Scope: every changed file matches `paths` (plus this task file). `git mv` keeps history for the moved tree.
- Manual check: `./ratatoskr --listen 127.0.0.1:0` logs the bound address. SIGINT logs `ratatoskr stopped` and exits 0.
- `git grep bbolt` outside docs: none. `admin/` and `pkg/` are gone.

**AC coverage**
- No `admin/`, no bbolt in `go.mod`: verified, and `go mod tidy` also dropped ff/v3, go-cmp, go-homedir and truststore.
- Placeholder on :8080: `defaultListen` is `127.0.0.1:8080` and binds loopback per docs/03 §8. Tests cover the handler
  (GET/HEAD/404/405 + headers), the real-listener `serve` lifecycle with ctx cancel, and flag parsing.
- Import guard: `go list -deps -json -test ./...`, skips synthetic `.test` mains and checks test imports, including
  external `_test` packages (their cut ImportPath is non-legacy). The `checked == 0` and empty-module-path sanity checks
  catch a guard that silently passes.
- Spec docs/03 §3: each package's fate matches the README table. Moving the whole `pkg/` tree (not only the five named
  packages) is justified in the Plan, because the guard would otherwise forbid `pkg/` → `internal/legacy/` imports.

**Findings**
- [ ] [minor] cmd/ratatoskr/main.go:55: `fs.SetOutput(io.Discard)` hides usage, so `./ratatoskr -h` prints
  `ERROR ratatoskr stopped error="parse flags: flag: help requested"` and exits 1. Handle `errors.Is(err, flag.ErrHelp)`
  in `run` by printing `fs.PrintDefaults()` (or keeping output on stderr only for help) and returning nil. Fine to fold into T-004, which
  replaces this CLI anyway; add it to T-004's AC so it isn't lost.
- [ ] [nit] cmd/ratatoskr/main.go:113: after the first signal cancels ctx, `stop()` only runs on return, so a second Ctrl-C
  during the 5 s drain is swallowed. Calling `stop()` right after `<-ctx.Done()` (pass it in, or re-arm in `main`)
  restores the default "second signal kills" behaviour.
- [ ] [nit] internal/legacy/importguard_test.go:96: the guard only sees the host GOOS/GOARCH and default build tags. A
  legacy import behind `//go:build windows` or a custom tag would slip through. Acceptable for now. If tagged files appear,
  run `go list` once per relevant `GOOS`/tags, or switch to a `go/parser` scan of all `.go` files.

**Out-of-scope items for the lead** (already listed by the implementer, confirmed still present)
- `gqlgen.yml` still points at `pkg/api/*` (inert; delete in T-007).
- `.gitignore:7` `/cmd/hetty/admin` is stale.
- `AGENTS.md:62` and `graphify-out/` still describe `admin/` and `pkg/`; refresh them after merge.

**Verdict: approve.** No blocker or major findings. Status stays `in_review` until the operator merges and sets `done`.
