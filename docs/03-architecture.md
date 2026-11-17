# 03 — Architecture

## 1. Big picture

```
 Browser / client ──► :8080 listener ──┬─► proxy engine ──► (rawhttp) ──► upstream / target
                                       │        │  intercept queue, match&replace, passive workflows
                                       │        ▼
                                       │   event bus ──► subscriptions (GraphQL WS) ──► UI
                                       │        │
                                       │        ▼
                                       │   store (SQLite: instance.db + projects/<id>/project.db)
                                       │
                                       └─► UI/API server (net/http): GraphQL, REST files, static SPA (go:embed)
```

Single process, single binary. The proxy front-end does **not** use `net/http` (fidelity, see ADR-0003);
the UI/API server does.

## 2. Repository layout (target after M0)

```
cmd/ratatoskr/                  main.go — flags, config, wiring via internal/app
internal/
  app/                      lifecycle: start/stop, project switching, dependency wiring (no globals)
  config/                   flags + env + config.toml (one struct, validated)
  ca/                       root CA, leaf issuance (SNI), cache, import/export
  rawhttp/                  lenient HTTP/1.x parser & writer preserving bytes; chunked; h2 client framing (M7)
  netx/                     dialer: upstream proxies, DNS rewrites, timeouts, TLS client config (SNI override)
  proxy/                    listeners, CONNECT, MITM, invisible mode, pipeline: capture → hooks → forward
  ws/                       websocket frame codec, permessage-deflate, stream recorder
  store/                    SQLite open/migrate, tx helpers, per-domain repositories (sub-packages)
    migrations/instance/    *.sql (goose-style numbered, embedded)
    migrations/project/     *.sql
  project/                  project registry, temp project, backups, restore, data deletion
  history/                  request/response records, WS streams/messages, sitemap derivation
  query/core|httpql|streamql  lexer, parser, AST, SQL compiler, evaluator, completions
  scope/  filterpreset/
  intercept/                hold queue (req/resp/ws), conditions, release semantics
  matchreplace/             rules, collections, engine, tester
  replay/                   collections, sessions, entries, sender, pipeline (last-byte sync)
  automate/                 sessions, payload sources, strategies, preprocessors, runner, extractors
  workflow/                 graph model, validation, executor, node registry
    nodes/                  one file per node family
  jsrt/                     JS runtime abstraction (goja), TS transpile (esbuild), SDK bindings
  findings/  exports/  files/  environment/
  assistant/                providers (anthropic, openai, google, openrouter, openai-compatible), chat store, tools
  plugins/                  package format, installer, backend host, frontend asset serving, community index
  events/                   typed in-process pub/sub (generic Bus[T]), backpressure policy
  tasks/                    background task manager (progress, cancel) for exports/backups/automate
  api/
    graphql/                schema/*.graphqls, gqlgen config, resolvers (thin: call services)
    rest/                   file upload/download, CA download, response preview, backup download
  auth/                     token auth, guest mode
  mcp/                      MCP server (M6)
  testutil/                 test servers (echo, chunked, slow, ws), golden helpers
web/                        Svelte SPA (see 06-frontend-conventions.md); build output web/dist embedded by internal/api/static
docs/ …
```

Rules: `internal/api/*` depends on services, never the reverse. Services depend on `store` repositories via small
interfaces defined **in the consumer package**. No package-level mutable globals. `context.Context` first arg everywhere.

## 3. What happens to Hetty code

| Hetty package | Fate | Notes |
|---|---|---|
| `admin/` (Next.js, React, MUI, Apollo) | **delete in T-002** | replaced by `web/` |
| `pkg/proxy/cert.go` | port → `internal/ca` | modernize (ECDSA, SNI cache, LRU) |
| `pkg/proxy/proxy.go`, `net.go`, `gzip.go`, `modify.go` | rewrite → `internal/proxy` + `rawhttp` | `httputil.ReverseProxy` removed (ADR-0003) |
| `pkg/proxy/intercept/*` | reference for `internal/intercept` | queue semantics reusable, filter replaced by HTTPQL |
| `pkg/filter/*` (own search syntax) | replace → `internal/query` | lexer/parser structure is a good model |
| `pkg/reqlog`, `pkg/sender` | rewrite → `history`, `replay` | |
| `pkg/scope` | port → `internal/scope` | add presets, globs |
| `pkg/proj` | rewrite → `internal/project` | multi-file DB model |
| `pkg/db/bolt` | delete | SQLite (ADR-0002) |
| `pkg/api` (gqlgen v0.14) | rewrite → `internal/api/graphql` | new schema, latest gqlgen, subscriptions |
| `pkg/chrome` (chromedp) | port → `internal/browser` (M7) | "open preconfigured Chromium" |
| `cmd/hetty` | rewrite → `cmd/ratatoskr` | |

Keep git history (fork), keep `LICENSE` copyright line, update module path `github.com/OWNER/ratatoskr`.

## 4. Proxy pipeline

```
accept conn → detect (CONNECT | absolute-form | origin-form→UI | invisible) 
  → [TLS MITM if CONNECT:443/any TLS ClientHello sniffed]
  → loop per request on conn:
      read raw request (rawhttp.ReadRequest, lenient, max header 1 MiB, body streamed with 64 MiB capture cap)
      record = history.NewPending(raw, conn meta)            // assigned ID, emitted requestAdded(pending)
      matchreplace.ApplyRequest(record)                       // sources=proxy
      intercept.HoldRequest(ctx, record) → forward|drop|edited
      passive workflows (On Intercept Request) — async unless they modify (then sync with timeout)
      forward via netx (keep-alive pool keyed by host:port:tls:sni:upstream)
      read raw response; matchreplace.ApplyResponse; intercept.HoldResponse; passive workflows (response)
      write response to client; history.Complete(record) → requestUpdated event
  → if Upgrade: websocket → ws.Recorder relays frames both ways (intercept/M&R hooks per message)
```
Bodies > capture cap are streamed through and stored truncated (`truncated=true`). Timeouts: dial 10 s, TLS 10 s,
response header 60 s (configurable).

## 5. Storage

- Driver: `modernc.org/sqlite` (pure Go). `PRAGMA journal_mode=WAL; synchronous=NORMAL; foreign_keys=ON; busy_timeout=5000`.
- One writer goroutine per DB (serialized writes via channel / `*sql.DB` with `SetMaxOpenConns(1)` for the write handle) + read pool.
- Batched inserts for proxy capture (flush every 50 ms or 200 rows).
- Migrations: embedded SQL files, numbered, forward-only, run at open; schema version table.

### 5.1 instance.db (main tables)
`settings(key, value_json)`, `projects(id, name, path, created_at, last_opened_at, is_default)`, `ca(cert_pem, key_pem_enc)`,
`workflows(id, kind, name, description, enabled, graph_json, shortcut, updated_at)`, `plugins(id, manifest_json, enabled, installed_at)`,
`plugin_storage(plugin_id, key, value)`, `hosted_files(id, name, size, sha256, path, created_at)`, `env_global(name, value, secret)`,
`ai_providers(id, kind, base_url, api_key_enc, models_json)`, `shortcuts(command_id, binding)`, `secrets_key` (file, not table).

### 5.2 project.db (main tables)
```
requests(id INTEGER PK, source TEXT, host TEXT, port INT, tls INT, method TEXT, path TEXT, query TEXT, ext TEXT,
         raw BLOB, body_offset INT, len INT, created_at INT, edited_from INT NULL, color TEXT NULL, note TEXT NULL,
         is_dropped INT DEFAULT 0, truncated INT DEFAULT 0, plugin_id TEXT NULL, automate_attack_id INT NULL)
responses(request_id PK/FK, code INT, raw BLOB, body_offset INT, len INT, roundtrip_ms INT, created_at INT, edited_from INT NULL)
request_headers(request_id, idx, name, name_lc, value)      -- indexed (name_lc), for header.* queries
response_headers(request_id, idx, name, name_lc, value)
ws_streams(id, request_id, host, port, path, tls, created_at, closed_at)
ws_messages(id, stream_id, direction, format, raw BLOB, len, created_at, edited_from)
scope_presets, filter_presets, ws_filter_presets
mr_collections, mr_rules(…, sort_order)
replay_collections, replay_sessions(kind http|ws|pipeline, …), replay_entries(session_id, request_id, created_at)
automate_sessions, automate_attacks, automate_results(attack_id, request_id, payloads_json, extractors_json, error)
findings(id, request_id, reporter, title, description_md, dedupe_key UNIQUE(reporter,dedupe_key), created_at)
exports(id, name, format, rows, size, path, status, created_at)
environments(id, name, selected), env_vars(env_id, name, value, secret)
assistant_sessions, assistant_messages
intercept_settings, ui_state(key, value_json)   -- per-project UI prefs (columns, pane sizes, scope selection)
```
Indexes: `requests(host)`, `(created_at)`, `(source)`, `(method)`, `responses(code)`, header `name_lc`.

## 6. API

- **GraphQL** (gqlgen, schema-first in `internal/api/graphql/schema/*.graphqls`, one file per domain).
  - Pagination: cursor connections (`first/after`, `last/before`) + `totalCount` (cached per query).
  - Subscriptions via `graphql-transport-ws` protocol on `/api/graphql/ws`: `requestAdded(query)`, `requestUpdated`,
    `interceptQueueChanged`, `wsMessageAdded(streamId)`, `automateProgress(attackId)`, `taskProgress`, `findingAdded`, `logLine`.
  - Errors: GraphQL errors with `extensions.code` (`BAD_QUERY`, `NOT_FOUND`, `CONFLICT`, `INTERNAL`) + `extensions.span` for query errors.
- **REST** (`/api/rest/…`): `POST /files` (multipart), `GET /files/{id}`, `GET /exports/{id}/download`, `GET /backups/{id}/download`,
  `GET /ca.crt|ca.der`, `GET /preview/{requestId}` (sandboxed response render), `GET /plugins/{id}/assets/*`.
- **Contract-first rule**: schema changes land in their own PR by the lead (see `07-agent-workflow.md`).
  Frontend types are generated (`bun run codegen`); backend resolvers generated by `go generate ./internal/api/graphql`.

## 7. Concurrency model
- One goroutine per client connection; per-request work bounded by context with deadline.
- `events.Bus[T]`: fan-out with per-subscriber buffered channel (cap 1024); slow subscribers get dropped + resync event (`stale=true`) instead of blocking the proxy.
- Passive workflows: worker pool (`GOMAXPROCS` workers), queue cap 10k, drop-oldest with metric.
- Automate runner: token-bucket rate limiter + worker pool; pause = stop dispatch, in-flight completes.

## 8. Security of the tool itself
- UI/API bound to 127.0.0.1 by default; warn loudly when bound to non-loopback without `--auth-token`.
- CSRF: API requires `Content-Type: application/json` + custom header `X-Ratatoskr: 1`; WS requires token in `connection_init` when auth on.
- Response preview served from a separate origin-like path with `Content-Security-Policy: sandbox`.
- Secrets (AI keys, secret env vars, CA key) encrypted at rest with XChaCha20-Poly1305 using a key file `secret.key` (0600).
- Shell node & plugins run with user privileges: show a one-time warning and a global "allow shell nodes" setting (default on, like Caido).

## 9. Observability
`log/slog` JSON/text handler; `slog.NewMultiHandler` fans out to stderr + in-memory ring buffer (for Logs page, 10k lines).
`/debug/pprof` (incl. `goroutineleak`) only when `--debug`.
