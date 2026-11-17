# 00 — Project Overview

> Codename: **Ratatoskr** (the Norse squirrel that carries — and tampers with — messages up and down Yggdrasil: the original man-in-the-middle). CLI alias: `rtsk`.
> Rename freely; the name appears only in `go.mod`, `cmd/`, `web/index.html`, and this doc set.

## 1. What we are building

Ratatoskr is a **free, open-source (MIT), single-binary web security proxy** for penetration testers and
bug-bounty hunters. It is a hard fork of [Hetty](https://github.com/dstotijn/hetty) (MIT, dormant since
v0.7.0 / 2022) and targets **feature parity with Caido** (as of Caido v0.58, Aug 2026), with no feature
locked behind a paid tier.

- Backend: **Go 1.27**, one static binary (`CGO_ENABLED=0`), embedded web UI.
- Frontend: **Svelte 5 SPA** (Vite, TypeScript), built and managed with **Bun**.
- Storage: **SQLite** (pure Go driver), one DB per project + one instance DB.
- API: **GraphQL** (queries, mutations, subscriptions over WebSocket) + small REST surface for file transfer.
- Client/server split like Caido: the UI talks to the backend only through the API, so a backend can run
  on a VPS and be used from a local browser.

## 2. Goals

1. Every feature listed in `docs/01-feature-spec.md` works, offline, with no account and no telemetry.
2. Raw-byte fidelity: what the user sends in Replay/Automate is exactly what goes on the wire.
3. Fast with large projects: 1M+ requests in history, table scrolling at 60 fps, filters < 300 ms on 100k rows.
4. Keyboard-first, dense, customizable UI (resizable panes, rebindable shortcuts, command palette).
5. Extensible: workflows (visual automation), JavaScript plugins (backend + frontend), MCP server for AI agents.

## 3. Non-goals (for v1.0)

- Cloud accounts, licensing, team billing, telemetry, auto-crawler/active scanner engine (a scanner can
  later be a plugin, as in Caido).
- A desktop wrapper. v1.0 ships as a CLI binary that serves the UI on `http://127.0.0.1:8080`.
  (A Tauri/Wails wrapper is a post-1.0 option, see ADR backlog.)
- Collaboration / multi-user real-time editing (post-1.0).

## 4. Legal & branding guardrails (MANDATORY for every agent)

Feature parity means **re-implementing behaviour**, not copying a product.

- **Never** copy Caido source, binaries, assets, icons, screenshots, CSS, documentation text, or its logo,
  colours, or name. Do not decompile or reverse-engineer Caido binaries (EULA).
- Caido's public docs are used only as a **behavioural reference**. Write all UI copy and docs in our own words.
- The **query language syntax is compatible with HTTPQL** (a public syntax; Caido's parser crate `httpql` is
  MIT-licensed). We implement our own parser in Go. In the UI we call it "HTTPQL-compatible query".
- Keep Hetty's original `LICENSE` copyright line and add ours below it. Every new file is MIT.
- Visual identity (colours, icons, typography, logo) is our own: see `docs/02-ui-layout.md §10`.

## 5. Glossary

| Term | Meaning |
|---|---|
| Instance | One running `ratatoskr` backend with its data directory. |
| Project | Isolated workspace (own DB file): history, replay, findings, scopes… |
| Request record | One proxied HTTP request + optional response, immutable once stored. Edits create new records. |
| Source | Which feature produced a request: `proxy`, `intercept`, `replay`, `automate`, `workflow`, `plugin`. |
| Scope preset | Named allow/deny host-glob lists. One preset is "active" per page. |
| Filter preset | Named, aliased HTTPQL query (referenced as `preset:"name"` / `preset:alias`). |
| Session (Replay) | A tab holding one editable request and its send history (entries). |
| Collection | Folder of Replay sessions. |
| Automate session | Template request + placeholders + payload config; each run produces an "attack" result table. |
| Workflow | Node graph. Types: **passive** (auto on traffic), **active** (manual on selected requests), **convert** (bytes → bytes). |
| Finding | A reported issue linked to a request (title, reporter, description in Markdown, dedupe key). |
| Hosted file | A file uploaded to the instance (wordlists etc.), used by Automate and plugins. |
| Environment | Named set of variables, interpolated as `{{name}}` in Replay/Automate; one Global + N custom, one selected. |

## 6. Document map

| File | Owner role | Purpose |
|---|---|---|
| `AGENTS.md` | lead | Rules every agent must follow. Read first, always. |
| `docs/00-overview.md` | lead | This file. |
| `docs/01-feature-spec.md` | lead | Feature-by-feature behaviour + acceptance criteria (the "what"). |
| `docs/02-ui-layout.md` | lead / frontend | App shell, page layouts, shortcuts, context menus, design tokens. |
| `docs/03-architecture.md` | lead / backend | Packages, data flow, storage schema, API, Hetty migration. |
| `docs/04-query-language.md` | backend | HTTPQL-compatible + StreamQL grammar, semantics, SQL compilation. |
| `docs/05-go-conventions.md` | backend | Go 1.27 idioms, libraries, testing, error handling. |
| `docs/06-frontend-conventions.md` | frontend | Svelte 5 runes, Bun, Vite, state, API client, components, tests. |
| `docs/07-agent-workflow.md` | lead | Multi-agent process: roles, worktrees, task lifecycle, review. |
| `docs/08-roadmap.md` | lead | Milestones M0–M7 and the task index. |
| `docs/tasks/` | everyone | One file per task (`T-###-slug.md`) + `TEMPLATE.md`. |
| `docs/adr/` | lead | Architecture Decision Records. Change a decision = new ADR. |
| `docs/PANDUAN.md` | human | Indonesian guide for the human operator. |
