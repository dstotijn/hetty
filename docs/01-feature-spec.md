# 01 — Feature Specification (parity target: Caido v0.58, Aug 2026)

Each feature has an ID (`F-<AREA>-<n>`). Tasks reference these IDs. "MUST" = required for parity.
Acceptance criteria (AC) are written so a QA agent can turn them into tests.
Layout details live in `docs/02-ui-layout.md`; query syntax in `docs/04-query-language.md`.

Paid-only features in Caido (Assistant, Automate preprocessors, workflow-in-automate, filtered export,
unlimited projects/plugins) are **free and unlimited** here.

---

## A. Core proxy & networking (`F-NET`)

### F-NET-1 Listener
- Default listener `127.0.0.1:8080` serves **both** the proxy and the UI/API (like Hetty/Caido):
  - absolute-form request (`GET http://…`) or `CONNECT` → proxy;
  - origin-form request with `Host` = listener address or `ratatoskr.localhost` → UI/API.
- Multiple listeners (address:port list) configurable in Settings → Network and via CLI `--listen` (repeatable).
- AC: `curl -x http://127.0.0.1:8080 http://example.test/` is recorded in history; `http://127.0.0.1:8080/` serves the UI.

### F-NET-2 HTTPS interception (MITM)
- Root CA generated on first start (ECDSA P-256, 10y), stored in instance data dir.
- Per-host leaf certificates issued on the fly from **SNI** (fallback: CONNECT host), cached (LRU, 1024).
- CA can be downloaded (PEM + DER) at `http://ratatoskr.localhost/ca.crt` and from Settings → Certificate;
  can be regenerated, imported (PEM cert+key) and exported.
- ALPN to client: `http/1.1` only (v1.0). Upstream: HTTP/1.1 by default.
- AC: browser with CA installed loads `https://example.com` with no warning; request appears with `tls=true`.

### F-NET-3 Raw fidelity
- The proxy parses HTTP/1.x itself (lenient parser); unknown/malformed-but-forwardable messages are forwarded byte-exact.
- Header order, case, duplicate headers, odd whitespace are preserved in storage and on the wire.
- Bodies stored raw (still compressed/chunked as on the wire); UI shows decoded view on demand.
- AC: golden tests with 30+ malformed/edge requests (duplicate CL, obs-fold, lowercase methods, absolute paths with fragments) round-trip byte-exact.

### F-NET-4 WebSocket proxying
- Upgrade handshakes are recorded as normal requests; subsequent frames are recorded as WS messages (see F-WS).
- Supports permessage-deflate (store decompressed payload + flag), continuation frames, ping/pong/close.

### F-NET-5 Upstream proxies
- Ordered list of rules: `{enabled, host glob, type http|socks5, address, auth}`; first match wins; "direct" fallback.

### F-NET-6 Invisible (transparent) proxying
- Listener option: treat origin-form requests as proxy traffic, routing by `Host` header / SNI; used for thick clients
  with DNS/hosts redirection. Optional "redirect to" host:port override per listener.

### F-NET-7 DNS rewrites
- Rules `{host glob → IP}` applied before dialing; plus custom resolver address.

### F-NET-8 Domain allowlist for UI access & local traffic
- Option to proxy (or bypass) traffic to localhost targets; list of hostnames allowed to reach the UI.

### F-NET-9 HTTP/2 upstream (beta flag)
- Replay and Replay Pipeline can send over HTTP/2 when enabled per session (h2 framing done by our engine).

---

## B. Workspace (`F-WSP`)

### F-WSP-1 Projects
- Create / rename / delete / duplicate / select project. Mark one as default (opened at start).
- A **temporary project** is created when none exists (deleted on exit unless saved).
- Project switch is instant for the UI (all stores reset + re-subscribe).
- Unlimited projects.
- AC: deleting a project removes its directory; switching projects never shows stale rows.

### F-WSP-2 Backups
- Create backup of a project (zip: `project.db` + `files/` + manifest JSON), list, download, restore as new project, delete.
- Backups run as background tasks with progress (see F-UI-9 Tasks).

### F-WSP-3 Data management
- Delete requests matching an HTTPQL query (with preview count). Vacuum DB. Show DB size per project.

### F-WSP-4 Instance auth (remote use)
- `--auth-token` or Settings: when set, UI/API require bearer token (cookie after login page). Guest mode = no token on localhost.

---

## C. Overview group

### F-SITE-1 Sitemap
- Tree derived from history: `scheme://host:port` → path segments → (leaf) request list. Updates live.
- Respects the active Scope preset. Node badges: request count; icon by type (folder/file/query).
- Select node → bottom table of requests under that node (same table component as HTTP History) + request/response panes.
- Shift-click expands a node fully. Context menu: delete node (deletes underlying requests after confirm), add host to scope, copy URL.
- AC: 10k requests over 200 hosts render the tree in < 500 ms; deleting a node removes its rows from history.

### F-SCOPE-1 Scopes
- Scope preset = `{name, allowlist[], denylist[]}`; entries are host globs (`*.example.com`, `api.*.com`), optional
  port and path prefix. Deny wins over allow. Empty allowlist = everything allowed.
- Create / edit / duplicate / delete / reorder presets. Each traffic page has a scope selector (none or one preset).
- "Add to scope" / "Exclude from scope" context actions everywhere a request appears.
- Intercept and workflows can be limited to the active scope.

### F-FILTER-1 Filter presets
- Preset = `{name, alias (lowercase-kebab), query}`. CRUD + import/export JSON.
- Referencable from any query as `preset:"Name"` or `preset:alias`. Cycles detected → error.
- Pages show quick-apply chips of presets next to the query bar.

---

## D. Proxy group

### F-INT-1 Intercept
- Global toggle in top bar: **Forwarding** (pass-through) ⇄ **Queuing** (hold). Separate toggles for
  Requests, Responses, WebSocket messages.
- Optional HTTPQL condition + scope preset deciding which traffic is held (others pass through).
- Queue tables (requests / responses / ws messages) with: id, host, method, path, time-held. Badge with queue count on sidebar item.
- Selected item opens editable raw editor. Buttons: **Forward**, **Drop**, Forward all, Drop all; undo edits.
- Actions: send to Replay / Automate, add to scope, run active workflow, copy as cURL.
- Edited requests are stored as a new record linked to the original (`edited_from`), both visible (F-HIST-4).
- When switching to Forwarding, all held items are released in FIFO order.
- AC: with queuing on, a browser request blocks until Forward; Drop returns `HTTP 502` (configurable) to client and marks record dropped.

### F-HIST-1 HTTP History table
- Virtualized table of all requests in the project; live insertion at top or bottom (sort-aware) via subscription.
- Columns (show/hide, reorder, resize, persisted): ID, Host, Method, Path, Query, Status, Length (resp), Ext,
  Roundtrip ms, Created at, TLS, Port, Source, Color, Comment/Note, Edited flag. Sort by any column.
- Query bar (HTTPQL-compatible, autocomplete, syntax highlight, error underline) + scope selector + preset chips.
- Row colour (8 colours + none) and note per row; set manually or by workflows.
- Multi-select (shift/ctrl). Bulk actions: send to Replay, send to Automate (first), run active workflow, export selection, delete, colour.
- AC: 100k rows scroll smoothly; filter result for simple host query on 100k rows < 300 ms.

### F-HIST-2 Request/response viewer
- Split view below (or beside, toggle) the table: request pane, response pane, each with modes:
  **Raw**, **Pretty** (formatted JSON/XML/HTML/JS by content-type), **Hex**, **Headers/Params table**, **Preview**
  (sandboxed iframe for HTML/images), plus "open response in browser" (served from `ratatoskr.localhost/preview/{id}`).
- Read-only in history; selection → context menu (send selection to Decoder-like convert workflow, copy, search).
- Auto-decoding of `Content-Encoding` (gzip, deflate, br, zstd) and chunked for display only.

### F-HIST-3 Context menu (rows & panes)
Send to Replay (Ctrl+R) · Send to Automate (Ctrl+M) · Run active workflow ▸ · Add host to scope / exclude ·
Create finding · Copy URL / as cURL / as raw / as fetch() / as Python requests · Set colour ▸ · Edit note ·
Open in browser · Filter by this host/path/status · Delete · Export.

### F-HIST-4 Modifications
- Records edited by Intercept or Match & Replace show an "edited" indicator; viewer can toggle **Original / Modified**.

### F-WS-1 WS History
- Top table: WebSocket **streams** (id, host, path, direction count, created, TLS, source).
- Bottom: messages of selected stream (direction ↑/↓, format text/binary, length, time) + message viewer (raw/pretty/hex).
- Query bar with **StreamQL** (see `04-query-language.md`). Send message to Replay (WS session).

### F-MR-1 Match & Replace
- Rules organised in **collections** (folders), ordered (drag to reorder), each `enabled` toggle.
- Rule = `{name, condition HTTPQL (optional), sources[] (proxy, replay, automate, workflow, plugin, intercept), section, operation, matcher, replacer}`.
- Sections & operations:
  - Request: First line (update raw), Method (update), Path (update raw), Query (update raw, update param, add param, remove param), Header (update raw, update value, add, remove), Body (update raw)
  - Response: First line (update raw), Status code (update), Header (update raw, update value, add, remove), Body (update raw)
  - WebSocket message: client→server, server→client (update raw)
- Matchers: **Full** (whole section / all items), **String** (exact), **Regex** (RE2 syntax — no look-around/backrefs; document this), **Name** (header/param name).
- Replacers: **Term** (string; regex capture groups `$1`, `${name}`), **Workflow** (convert workflow applied to the match).
- Body encoding awareness: decode `Content-Encoding` before matching, re-encode after, fix `Content-Length`.
- **Tester panel**: paste/pick a request/response → shows result diff live.
- Also applies in Replay when the session toggle "Apply Match & Replace" is on.
- Import/export JSON.
- AC: rule `Request Header / add / User-Agent: x` adds header on every proxied request; record marked edited; tester shows diff.

---

## E. Testing group

### F-REP-1 Replay sessions & collections
- Left panel: collections tree (Default collection always exists) with sessions; drag & drop sessions between
  collections; reorder; rename; duplicate; delete; search sessions by name/host/path.
- Top: tab strip of open sessions (close, close others, close to the right, rename by double-click, tag colour).
- Session types: **HTTP**, **WebSocket**, **Pipeline**.
- New session: from scratch (URL dialog with connection info) or sent from any request (Ctrl+R).

### F-REP-2 HTTP session editor
- Connection bar: scheme/host/port/TLS/SNI override, HTTP/1.1 vs HTTP/2 (beta), "Apply Match & Replace" toggle, environment selector.
- Request editor (CodeMirror, raw HTTP, syntax highlight, `{{var}}` highlighting & hover value).
- **Send** (Ctrl+Enter), **Cancel** in-flight. Options: follow redirects (off by default), update Content-Length automatically (on), timeout.
- Response pane with same modes as F-HIST-2 + roundtrip time, size, status.
- Entry history per session: `<` `>` navigation and history dropdown; each send = one stored entry (request + response).
- AC: sending `GET / HTTP/1.1` with a deliberately wrong `Content-Length` and auto-update off sends that exact byte stream.

### F-REP-3 WebSocket session
- Connect using the upgrade request; message composer (text/binary), message log both directions, reconnect, close.

### F-REP-4 Pipeline session (race conditions)
- List of requests (add from existing sessions or history). Strategies: **Sequential** (ordered, on one or many connections)
  and **Last-byte synchronization** (all requests written except final byte, then final bytes flushed together).
- Requests table shows per-request status; Results tab shows all responses side by side.
- AC: integration test with local server measuring arrival spread < 5 ms for 20 requests with last-byte sync.

### F-REP-5 Replay workflows
- From the request editor context menu: run convert workflows on selection (replace selection with output).

### F-AUTO-1 Automate session
- Same left panel structure as Replay (sessions list). Editor with request template and **placeholders**:
  select text → "Add placeholder" (+), remove (–), clear all. Placeholder default value = original text.
- Strategies:
  - **Sequential** (Sniper): one placeholder at a time, others keep original.
  - **All** (Battering ram): same payload in all placeholders.
  - **Parallel** (Pitchfork): one payload set per placeholder, iterate in lockstep (sets must be same length or stop at shortest — show warning).
  - **Matrix** (Cluster bomb): cartesian product.
- Payload types per set: **Simple list** (paste/edit), **Hosted file** (F-FILE), **Numbers** (from, to, step, padding, hex/dec), **Null payloads** (repeat N times / forever).
- **Preprocessors** (ordered, per set): Prefix, Suffix, URL-encode (charset field, non-ASCII toggle), Workflow (convert workflow), Base64, Hash.
- **Settings**: workers (concurrency), delay between requests (ms, jitter), retries on error + backoff, timeout,
  update Content-Length, follow redirects, HTTP version, apply Match & Replace, rate-limit (req/s).
- **Run** creates a paired results tab; Pause / Resume / Cancel; progress bar; ETA.
- Results table: payload columns, status, length, roundtrip, error, plus **extractors** (custom columns via regex group or JSON path on response) — sortable and HTTPQL-filterable (`row`, `resp.*`).
- Results rows open request/response viewer; send to Replay; create finding.
- AC: Matrix with sets of 3 and 4 values produces exactly 12 requests with correct substitutions; Pause halts dispatch within 1 s.

### F-WF-1 Workflows
- Page with three tabs: **Passive**, **Active**, **Convert**. List of workflows with enable toggle (passive), name, description,
  "global shortcut" (active/convert can be bound to a shortcut), import/export JSON, duplicate, delete.
- Visual editor (node graph): palette of nodes on the left, canvas centre, node inspector on the right, bottom run/test panel.
- Node inputs are **constant** values or **references** `$node_alias.output_name`; string **interpolation** `{{ $node.out }}` inside text inputs.
- Node aliases: `[a-z0-9_-]+`, unique per workflow.
- Execution: DAG from the start node; If/Else branches; each node gets typed inputs, returns typed outputs; errors stop the run and log.
- Node catalogue (MUST):

| Node | Passive | Active | Convert |
|---|:-:|:-:|:-:|
| On Intercept Request / On Intercept Response (start) | ✓ | | |
| Active Start / Active End | | ✓ | |
| Convert Start / Convert End | | | ✓ |
| Passive End | ✓ | | |
| If/Else (boolean input) | ✓ | ✓ | ✓ |
| If/Else JavaScript | ✓ | ✓ | ✓ |
| JavaScript | ✓ | ✓ | ✓ |
| Shell (choose shell, stdin = input, env vars for req/resp) | ✓ | ✓ | ✓ |
| In Scope (scope preset) | ✓ | ✓ | |
| Matches HTTPQL | ✓ | ✓ | |
| Set Color (row colour) | ✓ | ✓ | |
| Check Finding / Create Finding | ✓ | ✓ | |
| Logging | ✓ | ✓ | |
| Send Request (send modified request, output response) | | ✓ | |
| Base64 / Hex / HTML / URL encode & decode | | | ✓ |
| Join, Trim, JSON Minify, JSON Prettify, JWT Decode | | | ✓ |
| Match & Replace (regex/string on input) | | | ✓ |
| MD5 / SHA1 / SHA2 (256/384/512) hash | | | ✓ |

- JavaScript nodes: TypeScript/JS source, `export async function run(input, sdk)`; SDK gives request/response objects,
  `sdk.console`, `sdk.findings`, `sdk.requests.send`, `sdk.env`, `sdk.scope`, encoders. Typed definitions shipped to the editor.
- Integrations: Active workflows appear in request context menus; Convert workflows appear in editor selection menus,
  as Match & Replace replacers, and as Automate preprocessors. Passive workflows run on proxied traffic (bounded worker pool, never block forwarding beyond a configurable timeout).
- Test panel: run against a chosen request (passive/active) or input text (convert) and show per-node outputs + logs.
- Community/default workflows bundled: Base64 encode/decode, URL encode/decode, JWT decode, MD5, "colour in-scope GET".
- AC: convert workflow `Convert Start → Base64 Encode → Convert End` on input `abc` returns `YWJj`; passive workflow `In Scope → Matches HTTPQL resp.code.eq:500 → Set Color red` colours matching rows only.

### F-AI-1 Assistant (bring-your-own-key)
- Settings → AI: providers **Anthropic, OpenAI, Google, OpenRouter, OpenAI-compatible base URL (Ollama, LiteLLM, vLLM)**; keys stored encrypted at rest (instance key file); model list fetched or typed.
- Assistant page: chat sessions list (left), chat (centre), attached context chips (requests, selections). Streaming responses, Markdown + code blocks with "copy" and "send to Replay".
- Quick actions from context menu: **Explain request**, **Suggest attack vectors**, **Generate CSRF PoC** (HTML), **Write payloads for selection**.
- Privacy notice before first send; per-session toggle "redact cookies/Authorization".
- Tool-use (phase 2): the model may call safe tools (search history, read request, create replay session, add finding) — every tool call shown and user-approvable.
- AC: with a mock provider, "Explain request" sends the raw request (redacted when toggle on) and streams the answer.

---

## F. Logging group

### F-SEARCH-1 Search
- Like HTTP History but spans **all sources** (proxy, replay, automate, workflow, plugin, intercept) with a `source:` filter
  and source-inclusion toggles in advanced options. Same table, viewer, context menu.

### F-FIND-1 Findings
- Left: reporters list (All, Manual, each workflow/plugin name) with counts. Centre: findings table (title, host, path, reporter, created, request id).
  Right/bottom: finding detail (Markdown description editable, linked request/response viewer).
- Create manually from any request. Dedupe key (`reporter + dedupe_key`) prevents duplicates. Delete, export (JSON/CSV/Markdown).
- HTTPQL filtering on linked request + finding fields (`finding.title`, `finding.reporter` — our extension).

### F-EXP-1 Exports
- From History/Search/Automate results: export **all** or **current rows** (current filter + scope) as **JSON** or **CSV**
  (CSV flattens response columns with `response_` prefix). Large exports run as background tasks.
- Exports page lists exports: name, format, rows, size, status, created; download, rename, delete.

---

## G. Workspace group (pages)

### F-ENV-1 Environment
- Global environment (instance-wide) + per-project environments. One "selected" environment per project.
- Variables `{name, value, kind: plain|secret}`; secrets masked in UI and never sent to AI or exports.
- Interpolation `{{name}}` in Replay, Automate, workflows' Send Request; selected env overrides Global.
- Workflows/plugins can set variables (e.g. refreshed JWT).

### F-FILE-1 Files (hosted files)
- Upload (drag & drop, multi), list (name, size, created), rename, delete, preview first N lines. Used by Automate and plugins.

### F-PLUG-1 Plugins
- Plugin package = zip with `manifest.json` + `backend/*.js` and/or `frontend/*.js` + `frontend/*.css`.
- Install from file, from URL, or from a **community index URL** (JSON list: name, author, version, downloads, description, repo).
- Installed tab: enable/disable, update, uninstall, view logs, settings. Store tab: filter Official/Community, sort, multi-select install/update.
- Backend plugins run in a sandboxed JS runtime with a typed SDK (requests, findings, scope, env, files, storage, HTTP send, events: onInterceptRequest/Response, onProjectChange).
- Frontend plugins get a typed SDK: register pages, sidebar items, commands, command-palette entries, context-menu items, footer slots, editor extensions, call own backend functions, storage, toasts.
- Plugin API is ours (not Caido-compatible); document it in `docs/plugins-sdk.md` when M6 starts.

---

## H. Cross-cutting UI (`F-UI`)

- **F-UI-1 Command palette** (Ctrl/⌘+K): go to page/plugin, run command, run workflow, switch project, toggle intercept, search sessions.
- **F-UI-2 Shortcuts**: all commands rebindable in Settings → Shortcuts; one shortcut may map to multiple contexts; conflicts warned.
- **F-UI-3 Resizable panes** everywhere; sizes persisted per page.
- **F-UI-4 Sidebar customization**: reorder pages & plugin pages, pin, hide items/categories, create categories, reset to defaults; collapse sidebar.
- **F-UI-5 Appearance**: dark/light/system theme, UI font size, editor font size, density, custom CSS (Settings tab), custom JS (Settings tab, runs in UI with SDK).
- **F-UI-6 Editors**: CodeMirror 6, HTTP syntax mode, search (Ctrl+F), fold JSON, wrap toggle, prettify, "show non-printable", large-message virtualization (> 5 MB → read-only chunked view).
- **F-UI-7 Toasts & notifications**, **F-UI-8 Onboarding checklist** (CA install, browser proxy, first project).
- **F-UI-9 Background tasks** indicator (exports, backups, automate runs, plugin installs) in the top bar.
- **F-UI-10 Logs viewer** (Settings → Logs) streaming backend logs; debug mode toggle.

## I. CLI (`F-CLI`)
`ratatoskr [--listen addr]… [--data-dir path] [--ui-listen addr] [--no-open] [--auth-token t] [--debug] [--allow-guests]`
`ratatoskr ca export|regenerate` · `ratatoskr project list|backup|restore` · `ratatoskr version`.
Env vars mirror flags (`RATATOSKR_LISTEN`, …). Config file `config.toml` in data dir.

## J. MCP server (`F-MCP`, M6)
Optional MCP endpoint (`/mcp`, streamable HTTP, token-protected) exposing tools: search history (HTTPQL), get request,
send request, create replay session, list/create findings, list scopes, set environment variable. Lets external AI agents drive Ratatoskr.
