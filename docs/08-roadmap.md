# 08 — Roadmap & Task Index

Milestones are sequential; tasks inside a milestone run in parallel where `depends_on` allows.
Only **M0 and M1** have full task files in `docs/tasks/`. The lead expands the next milestone into task files
(from `TEMPLATE.md`) when the current one is ~70% done — specs evolve, so don't pre-write everything.

Legend — role: L lead · B backend · F frontend · Q qa. Status lives in each task file.

## M0 — Foundation (repo runs, new skeleton, nothing user-facing yet)

| ID | Title | Role | Depends on |
|---|---|---|---|
| T-001 | Repo bootstrap: module path, Go 1.27, `internal/` layout, Makefile, golangci, CI (Linux/macOS), LICENSE attribution | L | — |
| T-002 | Remove Hetty `admin/` + bbolt; keep app compiling with stubbed services | B | T-001 |
| T-003 | Svelte 5 + Vite + Bun scaffold, tokens, app shell (TopBar, Sidebar, StatusBar, router, command palette skeleton), embed into Go | F | T-001 |
| T-004 | `internal/config` + CLI flags/env/TOML + `cmd/ratatoskr` wiring + slog multi-handler + `/debug` | B | T-002 |
| T-005 | `internal/store`: SQLite open, WAL, migrations runner, write serializer, test helpers; instance + project schemas v1 | B | T-002 |
| T-006 | `internal/project`: registry, temp project, switch, delete (contract added in T-007) | B | T-005, T-007 |
| T-007 | GraphQL server scaffold: gqlgen latest, schema dir, subscriptions transport, error codes, auth/CSRF middleware, codegen on web | L | T-003, T-004 |
| T-008 | `internal/ca`: root CA gen/import/export, SNI leaf cache, REST CA download | B | T-004 (core), T-007 (REST) |
| T-009 | Workspace page (projects tab) + project switcher in top bar | F | T-003, T-006, T-007 |

## M1 — Proxy core, History, Query language, Scopes, Sitemap

| ID | Title | Role | Depends on |
|---|---|---|---|
| T-010 | `internal/rawhttp`: lenient reader/writer, chunked, keep-alive, golden + fuzz tests | B | T-001 |
| T-011 | `internal/proxy` + `netx`: listener multiplexing (proxy vs UI), CONNECT, MITM, forward via rawhttp, capture to store (batched), events | B | T-005, T-008, T-010 |
| T-012 | Shared `TrafficTable` (virtualized) + `MessageViewer` (CodeMirror, Raw/Pretty/Hex/Table/Preview) | F | T-003 |
| T-013 | Contract: history queries/subscriptions, completions, row colour/note mutations | L | T-007 |
| T-014 | `internal/query`: HTTPQL lexer/parser/AST/errors with spans + conformance corpus | B | T-001 |
| T-015 | HTTPQL SQL compiler + SQLite functions + evaluator + benchmarks | B | T-014, T-005 |
| T-016 | History service + resolvers (pagination, filtering, live subscription filtering) | B | T-011, T-013, T-015 |
| T-017 | HTTP History page (query bar w/ HTTPQL CodeMirror mode + autocomplete, table, viewers, context menu v1) | F | T-012, T-013 |
| T-018 | Scopes + filter presets: backend services, `preset:` expansion, resolvers | B | T-015 |
| T-019 | Scopes & Filters pages + scope selector component + preset chips | F | T-017, T-018 |
| T-020 | Sitemap: backend tree derivation + subscription; Sitemap page | B+F | T-016 |
| T-021 | E2E: proxy a local target through Chromium, verify history/filter/scope/sitemap (Playwright) | Q | T-017, T-019, T-020 |
| T-022 | Command registry, shortcuts manager (rebindable, persisted), context menu builder | F | T-003 |

## M2 — Intercept, Replay, Match & Replace, Environment

| ID | Title | Role |
|---|---|---|
| T-030 | Contract: intercept, replay, M&R, environment | L |
| T-031 | Intercept queue engine (req/resp), conditions, release semantics, drop response | B |
| T-032 | Intercept page + top-bar toggle + sidebar badge | F |
| T-033 | Replay backend: collections, sessions, entries, sender (rawhttp), cancel, redirects, env interpolation | B |
| T-034 | Replay page: session sidebar, tab strip, connection bar, editor with `{{var}}` support, history nav | F |
| T-035 | Match & Replace engine (sections/operations/matchers/term replacer), body encodings, tester | B |
| T-036 | Match & Replace page (collections, rule editor, tester diff) | F |
| T-037 | Environment backend + page + status-bar switcher | B+F |
| T-038 | Original/Modified viewer toggle + edited markers in history | B+F |
| T-039 | E2E for M2 | Q |

## M3 — Automate, Files, Search, Findings, Exports, Tasks

| ID | Title | Role |
|---|---|---|
| T-040 | Contract: automate, files, search, findings, exports, background tasks | L |
| T-041 | Hosted files (REST upload/download, service, page) | B+F |
| T-042 | Automate engine: payload sources (iter.Seq), strategies, preprocessors (non-workflow), settings, rate limit, pause/resume | B |
| T-043 | Automate extractors + results storage + HTTPQL on results | B |
| T-044 | Automate pages (session editor with placeholders, payload panel, results tab) | F |
| T-045 | Background task manager + top-bar indicator | B+F |
| T-046 | Search page (all sources) | B+F |
| T-047 | Findings backend + page (reporters, markdown, dedupe) | B+F |
| T-048 | Exports (JSON/CSV streaming, all/current rows) + page | B+F |
| T-049 | E2E for M3 | Q |

## M4 — WebSockets

| ID | Title | Role |
|---|---|---|
| T-050 | WS frame codec + permessage-deflate + recorder in proxy | B |
| T-051 | StreamQL (reuse core) + WS history service | B |
| T-052 | WS History page | F |
| T-053 | Intercept WS messages + M&R for WS | B+F |
| T-054 | Replay WS session | B+F |
| T-055 | E2E for M4 | Q |

## M5 — Workflows

| ID | Title | Role |
|---|---|---|
| T-060 | Contract: workflows (graph model, run/test API, logs) | L |
| T-061 | `jsrt` (goja pool, limits, esbuild TS transform, SDK bindings) | B |
| T-062 | Workflow engine (validation, DAG executor, references, interpolation, logging) | B |
| T-063 | Node families: start/end, control, JS, shell, convert encoders/hashes/JSON/JWT/M&R | B |
| T-064 | HTTP-aware nodes: In Scope, Matches HTTPQL, Set Color, Findings, Send Request; passive pipeline integration | B |
| T-065 | Workflow editor UI (@xyflow/svelte, palette, inspector, test panel) | F |
| T-066 | Integrations: active workflows in context menus, convert in editors, M&R workflow replacer, Automate workflow preprocessor, shortcuts | B+F |
| T-067 | Bundled default workflows + import/export | B |
| T-068 | E2E for M5 | Q |

## M6 — Plugins, Assistant, MCP

| ID | Title | Role |
|---|---|---|
| T-070 | Plugin package format + manifest schema + installer + community index | L+B |
| T-071 | Backend plugin host (jsrt) + SDK + events | B |
| T-072 | Frontend plugin host + SDK (pages, sidebar, commands, menus, slots, editor ext) | F |
| T-073 | Plugins page (Store/Installed) | F |
| T-074 | `docs/plugins-sdk.md` + example plugins (hello-world, autorize-like) | L |
| T-075 | AI providers + encrypted key storage + Settings → AI | B+F |
| T-076 | Assistant page + quick actions + redaction | B+F |
| T-077 | MCP server | B |
| T-078 | E2E for M6 | Q |

## M7 — Networking & polish → v1.0

| ID | Title | Role |
|---|---|---|
| T-080 | Upstream proxies (HTTP/SOCKS5), DNS rewrites, multiple listeners, invisible proxy, allowlist | B+F |
| T-081 | HTTP/2 client for Replay + Pipeline (last-byte sync for h1 and h2) | B |
| T-082 | Replay Pipeline UI | F |
| T-083 | Backups/restore, data deletion, vacuum | B+F |
| T-084 | Instance auth token + guest mode + remote access hardening | B+F |
| T-085 | Settings pages complete (appearance, custom CSS/JS, logs viewer, editor) | F |
| T-086 | Sidebar customization + onboarding checklist | F |
| T-087 | Preconfigured Chromium launcher (port of Hetty `pkg/chrome`) | B |
| T-088 | Packaging: goreleaser, Docker image, checksums, SBOM; docs site | L |
| T-089 | Performance pass (1M rows), security review, release 1.0 | L+Q |
