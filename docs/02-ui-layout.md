# 02 — UI Layout & Interaction Spec

The **information architecture** mirrors Caido's (users switching from Caido should find everything where
they expect it). The **visual identity** is our own (§10). Wireframes are ASCII; proportions are defaults —
every divider is draggable and persisted.

## 1. App shell

```
┌──────────────────────────────────────────────────────────────────────────────────────┐
│ ◆ Ratatoskr  [Project ▾ acme-pentest]   ⌘K Search…        [|| Queuing ▾]  ⟳2  🔔  👤 ▾ │  ← Top bar (40px)
├───────────────┬──────────────────────────────────────────────────────────────────────┤
│ OVERVIEW      │                                                                      │
│  Sitemap      │                                                                      │
│  Scopes       │                                                                      │
│  Filters      │                         <active page>                                │
│ PROXY         │                                                                      │
│  Intercept 3  │                                                                      │
│  HTTP History │                                                                      │
│  WS History   │                                                                      │
│  Match & Repl │                                                                      │
│ TESTING       │                                                                      │
│  Replay       │                                                                      │
│  Automate     │                                                                      │
│  Workflows    │                                                                      │
│  Assistant    │                                                                      │
│ LOGGING       │                                                                      │
│  Search       │                                                                      │
│  Findings     │                                                                      │
│  Exports      │                                                                      │
│ WORKSPACE     │                                                                      │
│  Environment  │                                                                      │
│  Files        │                                                                      │
│  Plugins      │                                                                      │
│  Workspace    │                                                                      │
│ PLUGINS       │  (plugin-registered pages appear here)                               │
│ ⚙ Customize  «│                                                                      │
├───────────────┴──────────────────────────────────────────────────────────────────────┤
│ ● 127.0.0.1:8080  proxy ok │ project 1.2 GB │ env: staging │ <plugin footer slots>   │  ← Status bar (24px)
└──────────────────────────────────────────────────────────────────────────────────────┘
```

- **Top bar**: logo (click → Sitemap), project switcher (search, create, manage), command palette trigger,
  intercept toggle (Forwarding ▶ / Queuing ‖, dropdown for Requests/Responses/WS toggles), background tasks
  indicator, notifications, user menu (Settings, Shortcuts, Logs, About, Theme).
- **Sidebar** (220px, collapsible to 48px icon rail with `«`): categories with uppercase labels; item badge for
  counts (Intercept queue, new Findings). "Customize" enters edit mode (drag, pin, hide, new category, reset).
- **Status bar**: listener status, project DB size, selected environment (click to switch), plugin slots.
- Routing: hash router `#/<page>[/<id>]` (e.g. `#/replay/session/42`). Back/forward work. Deep links survive reload.

## 2. Shared building blocks (implement once, reuse everywhere)

| Component | Used by | Notes |
|---|---|---|
| `QueryBar` | History, Search, Sitemap, Intercept, WS, Automate results, Findings | HTTPQL/StreamQL editor (single-line CodeMirror), autocomplete fields/operators/presets, error squiggle, history ↑↓, Enter applies. Left: scope selector; right: preset chips + "Advanced" popover. |
| `TrafficTable` | History, Search, Sitemap, Automate results, Intercept queue | Virtualized, sortable, column config menu, multi-select, row colour stripe, live updates, keyboard nav (↑↓, Enter opens, Space select). |
| `MessageViewer` | everywhere a req/resp is shown | Tabs Raw · Pretty · Hex · Table · Preview; toolbar: wrap, search, copy, prettify, original/modified switch, "open in browser". Read-only or editable. |
| `SplitPane` | all pages | Horizontal/vertical, persisted sizes, double-click divider resets. |
| `ContextMenu` | rows, editors, tabs, tree nodes | Built from a command registry (same commands as palette + shortcuts + plugins). |
| `SessionSidebar` | Replay, Automate, Assistant | Collections tree + search + drag/drop. |
| `TabStrip` | Replay, Automate | Rename, close variants, colour tag, overflow menu, drag reorder. |
| `EmptyState` | all | Short explanation + primary action. |

## 3. Page layouts

### 3.1 HTTP History / Search
```
┌ [Scope ▾] [ req.host.cont:"api" AND resp.code.gte:400        ] [preset chips] [⚙]┐
├───────────────────────────────────────────────────────────────────────────────────┤
│ ID │ Host            │ Method │ Path            │ Status │ Len   │ RTT │ Created  │
│ 812│ api.acme.test   │ POST   │ /v2/login       │ 401    │ 312   │ 88  │ 10:21:03 │
│ …  (virtualized)                                                                  │
├──────────────────────────────────────────┬────────────────────────────────────────┤
│ Request  [Raw|Pretty|Hex|Table]          │ Response [Raw|Pretty|Hex|Table|Preview]│
│ POST /v2/login HTTP/1.1                  │ HTTP/1.1 401 Unauthorized              │
│ …                                        │ …                                      │
└──────────────────────────────────────────┴────────────────────────────────────────┘
```
Default split: table 55% / viewers 45%; viewers side-by-side. Layout toggle: viewers right of table.
Search adds a **Source** column and source toggles in ⚙ Advanced.

### 3.2 Sitemap
Left 30%: tree (host → path). Right: `TrafficTable` (top) + `MessageViewer` pair (bottom). Query bar on top applies to tree too.

### 3.3 Scopes / Filters
Two-column CRUD: left list of presets (+ New), right editor form (name, alias for filters, allow/deny lists as editable chip lists or HTTPQL editor with live "matches N requests" counter).

### 3.4 Intercept
```
┌ [Requests ‖] [Responses ▶] [WS ▶]   condition: [ HTTPQL…        ] [Scope ▾]    ┐
├──────────────────────────────┬────────────────────────────────────────────────────┤
│ Queue (Requests | Responses | WS) │  Editable editor of selected item              │
│ #  host       method path      │  POST /api HTTP/1.1                               │
│ 12 acme.test  GET    /         │  …                                                │
│ 13 …                           │                                                   │
│                                │  [Forward ⌃⏎] [Drop ⌃⌫] [Forward all] [Undo]     │
└──────────────────────────────┴────────────────────────────────────────────────────┘
```

### 3.5 WS History
Top: streams table. Bottom-left: messages table (direction arrow, format, length, time). Bottom-right: message viewer. StreamQL query bar.

### 3.6 Match & Replace
```
┌ Collections & rules (30%)      │ Rule editor (40%)                │ Tester (30%)               ┐
│ ▾ Headers                      │ Name [ Add bug bounty header ]   │ Input  [pick request ▾]    │
│   ☑ Add X-Bug-Bounty           │ Condition [ HTTPQL… ]            │ ┌ before ┐ ┌ after ┐      │
│   ☐ Strip CSP                  │ Sources [proxy ☑ replay ☑ …]     │ │        │ │       │      │
│ ▾ Auth                         │ Section [Request Header ▾]       │ diff highlighted           │
│ + New rule  + New collection   │ Operation [Add ▾]                │                            │
│                                │ Matcher [Full|String|Regex|Name] │                            │
│                                │ Replacer [Term|Workflow] [value] │                            │
└────────────────────────────────┴──────────────────────────────────┴────────────────────────────┘
```

### 3.7 Replay
```
┌ Sessions (240px)   │ [tab: login ×] [tab: idor-user ×] [ + ]                                  ┐
│ 🔍 search          ├────────────────────────────────────────────────────────────────────────┤
│ ▾ Default          │ https ▾ │ api.acme.test │ 443 │ HTTP/1.1 ▾ │ M&R ☐ │ env: staging ▾ │ ⚙   │
│   login            │ [Send ⌃⏎] [Cancel]   ‹ 3/7 ›  history ▾                                │
│   idor-user        ├─────────────────────────────────┬──────────────────────────────────────┤
│ ▾ Gadgets          │ Request (editable)              │ Response  200 · 1.2 KB · 143 ms      │
│   ssrf-probe       │ GET /users/{{uid}} HTTP/1.1     │ HTTP/1.1 200 OK                      │
│ + Session  + Coll. │ Host: api.acme.test             │ …                                    │
└────────────────────┴─────────────────────────────────┴──────────────────────────────────────┘
```
Pipeline session replaces editor area with: requests table (order, status) + strategy selector + Results tab.
WS session: connection header, message composer at bottom, message log.

### 3.8 Automate
Same shell as Replay. Session body:
```
│ Request template (placeholders highlighted §…§)   │ Payload panel                             │
│ POST /login HTTP/1.1                              │ Strategy [Matrix ▾]                       │
│ user=§admin§&pass=§x§                             │ ▸ Set 1 (placeholder 1)  Type [Hosted ▾] │
│ [+ Add placeholder] [– Remove] [Clear]            │   Preprocessors [URL-encode] [+]          │
│                                                   │ ▸ Set 2 …                                 │
│                                                   │ Settings (workers, delay, retries, …)    │
│                                                   │ [Run ▶]                                   │
```
Run opens a paired results tab: progress bar + Pause/Resume/Cancel, `TrafficTable` with payload & extractor columns, viewers.

### 3.9 Workflows
Tabs Passive | Active | Convert. List view (table: enabled, name, description, shortcut). Editor:
```
│ Nodes palette (200px) │ Canvas (pan/zoom, minimap)            │ Inspector (280px)          │
│ ▸ Start/End           │   [On Intercept Request]               │ Alias [ check-500 ]        │
│ ▸ Control             │            │                           │ Inputs: query = [ … ]      │
│ ▸ Code                │   [Matches HTTPQL] ─true─▶ [Set Color] │   (●const ○ref)            │
│ ▸ Encoding / Hashing  │            └false─▶ [Passive End]      │ Outputs: matched: Boolean  │
│ ▸ Findings / Misc     │                                        │                            │
├───────────────────────┴────────────────────────────────────────┴────────────────────────────┤
│ Test: [pick request ▾ | input text]  [Save] [Save & Run]   per-node outputs & logs          │
```
Node colour by category: Start/End, Control, Code, Convert, Findings/Misc (our palette, §10).

### 3.10 Assistant
Left: chat sessions; centre: messages (streaming), context chips above input, model selector; right (collapsible): attached request viewer.

### 3.11 Findings
Left: reporters with counts; centre: findings table; right: detail (Markdown editor/preview + linked request/response).

### 3.12 Exports, Files, Plugins, Workspace, Environment
Simple table + toolbar pages. Plugins: tabs **Store** / **Installed**. Workspace: tabs **Projects** / **Backups**. Environment: left env list (Global pinned), right variables grid (name, value, secret toggle).

### 3.13 Settings (full page, from user menu)
Tabs: General · Appearance · Shortcuts · Network (listeners, upstream, DNS, invisible, allowlist) · Certificate · AI · Editor · Custom CSS · Custom JS · Logs · About.

## 4. Default shortcuts (all rebindable)

| Action | Win/Linux | macOS |
|---|---|---|
| Command palette | Ctrl+K | ⌘K |
| Toggle sidebar | Ctrl+B | ⌘B |
| Send to Replay | Ctrl+R | ⌘R |
| Send to Automate | Ctrl+M | ⌘M |
| Send (Replay) | Ctrl+Enter | ⌘⏎ |
| Forward (Intercept, item selected) | Ctrl+Enter | ⌘⏎ |
| Drop (Intercept, item selected) | Ctrl+Backspace | ⌘⌫ |
| Toggle intercept | Ctrl+I | ⌘I |
| Next / prev tab | Ctrl+] / Ctrl+[ | ⌘] / ⌘[ |
| Close tab | Ctrl+W (in-app) | ⌘W |
| Focus query bar | `/` (outside inputs) | `/` |
| Find in editor | Ctrl+F | ⌘F |
| Delete selection | Delete | ⌫ |
| Go to page | `g` then letter (h history, r replay, a automate, i intercept, s search, f findings, w workflows) | same |

Browser-reserved combos (Ctrl+W/T/N) only work when the UI runs as an installed PWA; document fallback (`Alt+W`).

## 5. Context-menu command registry
Every action is a **command** `{id, title, icon, when(ctx), run(ctx)}` registered centrally; menus, palette,
shortcuts and plugins all consume the same registry. Ids are namespaced: `history.sendToReplay`, `editor.convert.<workflowId>`…

## 6. Live data rules
- Tables subscribe to `requestAdded/requestUpdated` events filtered server-side by the active query+scope.
- New rows do not steal selection or scroll position; a "↓ 23 new" pill appears when the user is scrolled away.

## 7. Empty states & onboarding
First run shows a checklist card on Sitemap: install CA (link), configure browser proxy (FoxyProxy hint, "Open preconfigured Chromium" if available), send first request, create project.

## 8. Accessibility
Full keyboard navigation, visible focus rings, ARIA roles on tables/tree/tabs, contrast ≥ 4.5:1 in both themes, no colour-only signals (row colour also shown as a stripe + tooltip).

## 9. Performance budgets
First load ≤ 400 KB JS gzip (code-split per page), page switch < 100 ms, table frame budget 16 ms, editor handles 5 MB messages.

## 10. Visual identity (ours, not Caido's)
- Palette: neutral slate greys; accent **teal `#14b8a6`** (primary), warn amber, danger rose, success emerald.
  Row colours: red, orange, yellow, green, teal, blue, purple, pink.
- Typography: UI `Inter` (bundled), code `JetBrains Mono` (bundled, OFL). Base size 13px, dense line-height 1.35.
- Icons: Lucide (`@lucide/svelte`). Logo: simple wave-through-gap mark (create in `web/src/assets/logo.svg`).
- Radius 6px, 1px borders, subtle shadows only on overlays. Dark theme default.
