# AGENTS.md — Ratatoskr

Ratatoskr is an MIT-licensed web security proxy (fork of Hetty) targeting feature parity with Caido.
Backend **Go 1.27** single binary · Frontend **Svelte 5 SPA** built with **Bun** · Storage **SQLite** · API **GraphQL**.
These rules apply to every AI agent (Claude Code, OpenCode, any model). They override your defaults.

## Read before doing anything
1. This file (always).
2. Your task file `docs/tasks/T-###-*.md` — it lists your role, allowed `paths`, dependencies and the spec sections to read.
3. Only the spec sections your task links. Doc map: `docs/00-overview.md §6`.
Do not start coding without a task file. If you were asked to do something with no task file, ask the operator for one or
(lead role only) create it from `docs/tasks/TEMPLATE.md`.

## Hard rules (violations = automatic `changes_requested`)
1. **Stay in your lane.** Modify only files matching your task's `paths`. Need more? Write it under "Out-of-scope change needed" in the task file and stop.
2. **Spec is truth.** If code/spec conflict or the spec is ambiguous, set `status: blocked`, write the question, stop. Never invent product behaviour silently.
3. **No Caido code, assets, text, colours, icons or names in the product.** Behaviour parity only (`docs/00-overview.md §4`).
4. **Generated code is never edited by hand**: `internal/api/graphql/generated/**`, `web/src/lib/api/generated/**`. Change the schema (lead) and run `make gen`.
5. **GraphQL schema files are owned by the lead** (`internal/api/graphql/schema/**`). Backend/frontend agents request changes via the task file.
6. **Byte fidelity**: never route proxy/Replay/Automate traffic through `net/http` client/server or `httputil` (ADR-0003).
7. **Never** weaken tests to make them pass, skip lint rules inline without a justification comment + task id, or commit secrets/test data with real targets.
8. **Verify APIs you are unsure of** (`go doc pkg.Symbol`, package README/types in `node_modules`). Do not guess signatures or invent packages. Add dependencies only from the approved lists (`docs/05`, `docs/06`) or via an ADR.
9. Run `make check` before setting `in_review`. It must be green.
10. Don't touch `main` directly. Work in your task branch/worktree. Conventional Commits with task scope: `feat(T-014): …`.
11. **This repo is a GitHub fork of `dstotijn/hetty`.** Never open PRs, issues or pushes against the upstream repo. Always target our repo explicitly:
    `gh pr create --repo OWNER/ratatoskr --base main …`. If `gh` would default to `dstotijn/hetty`, stop and tell the operator.

## Tech quick reference
- Go: `go 1.27`, `CGO_ENABLED=0`, modern idioms required — see `docs/05-go-conventions.md` (generic methods, `encoding/json/v2`, iterators, `errors.AsType`, `new(expr)`, `sync.WaitGroup.Go`, `os.Root`, `testing/synctest`, `slog`).
- Frontend: Svelte 5 runes only (`$state/$derived/$props`, snippets, callback props), Vite, Tailwind v4, bits-ui, CodeMirror 6, urql + codegen — see `docs/06-frontend-conventions.md`. **Bun only** (`bun install`, `bun run`, `bunx --bun`), never npm/pnpm/yarn.
- Layout: `docs/03-architecture.md §2`. Services in `internal/<domain>`, API layer thin.

## Commands
```
make gen        # gqlgen + web codegen
make check      # go fix diff, vet, golangci-lint, go test -race, web check + unit tests
make test-int   # Go integration tests
make build      # web build + go build → ./ratatoskr
make run        # run with ./.dev-data on 127.0.0.1:8080
cd web && bun run dev        # UI dev server, proxies /api to :8080
cd web && bun run test:e2e   # Playwright (builds binary)
```
Parallel worktrees: run the app with `--data-dir ./.dev-data --listen 127.0.0.1:0` to avoid port clashes.

## Task lifecycle (summary of `docs/07-agent-workflow.md`)
`todo → in_progress → in_review → (changes_requested ↔ in_review) → done`, or `blocked`.
On start: set status/owner, write a ≤15-line plan in the task file. On finish: fill **Handoff** (what changed, how to verify,
limitations, follow-ups). Reviewers write findings in **Review** with severity and file:line; they never edit code.

## Definition of Done
Acceptance criteria covered by tests · `make check` green · E2E/integration added for user-visible behaviour ·
no untracked TODOs (`TODO(T-###)`) · docs updated · Handoff complete · reviewer approved.

## When in doubt
Prefer the simplest design that satisfies the spec and the performance budgets. Ask (blocked) rather than guess.

## Knowledge graph (graphify)
A code + docs knowledge graph lives in `graphify-out/` (`graph.json`, `GRAPH_REPORT.md`, `graph.html`).
- For architecture or "how does X connect to Y" questions, read `graphify-out/GRAPH_REPORT.md` first, or run `graphify query "<question>"`, before grepping raw files.
- It does not replace the spec: `docs/` stays the source of truth. Treat INFERRED/AMBIGUOUS edges as hints, not facts.
- After large changes, refresh with `graphify update .` (code only, no LLM) or `/graphify . --update` (includes docs).
- The graph still contains legacy Hetty code (`admin/`, `pkg/`) until T-002 removes it.
