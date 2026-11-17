# 06 — Frontend Conventions (Svelte 5 SPA + Bun)

## 1. Stack (versions as of Nov 2026 — always `bun add <pkg>@latest` and record in package.json)

| Concern | Choice |
|---|---|
| Runtime / package manager / scripts | **Bun** (`bun install`, `bun run …`, `bunx --bun vite`) — lockfile `bun.lock` committed. No npm/pnpm/yarn. |
| Framework | **Svelte 5** (runes) — plain Vite SPA, **no SvelteKit** (ADR-0005) |
| Bundler / dev server | Vite (+ `@sveltejs/vite-plugin-svelte`) |
| Language | TypeScript strict (`"strict": true, "noUncheckedIndexedAccess": true`). Use the TS version that `svelte-check` supports (check its peerDependencies before upgrading major versions). |
| Styling | Tailwind CSS v4 (CSS-first config, design tokens as CSS variables in `src/styles/tokens.css`) |
| Headless UI primitives | `bits-ui` (menus, dialogs, popovers, tabs, tooltips, context menu, combobox) |
| Resizable panes | `paneforge` |
| Code editors | CodeMirror 6 (`@codemirror/*`) with our own HTTP language + HTTPQL/StreamQL modes |
| Virtualized tables | `@tanstack/svelte-virtual` (or in-house `VirtualList` if it lags Svelte 5 — decide in T-012) |
| Node graph (workflows) | `@xyflow/svelte` |
| Icons | `@lucide/svelte` |
| GraphQL client | `@urql/svelte` + `graphql-ws`; types via `@graphql-codegen/cli` client preset (`bun run codegen`) |
| Markdown | `marked` + `DOMPurify` (findings, assistant) |
| Tests | `vitest` + `@testing-library/svelte` (unit/component), `@playwright/test` (E2E against a real `ratatoskr` binary) |
| Lint/format | `eslint` + `eslint-plugin-svelte` + `typescript-eslint`, `prettier` + `prettier-plugin-svelte` + `prettier-plugin-tailwindcss` |

## 2. Directory layout (`web/`)

```
web/
  index.html
  package.json  bun.lock  vite.config.ts  tsconfig.json  codegen.ts  eslint.config.js
  src/
    main.ts                 mount(App)
    App.svelte              shell: TopBar, Sidebar, <RouterView/>, StatusBar, CommandPalette, Toasts
    lib/
      router/               tiny hash router: routes table, `navigate()`, `route` rune state, lazy page imports
      api/                  urql client, ws client, generated/ (codegen output, never edited by hand), rest.ts
      commands/             command registry, shortcuts manager, context-menu builder
      stores/               *.svelte.ts global state classes (project, settings, intercept, environment, tasks, toasts)
      query/                HTTPQL/StreamQL CodeMirror language, autocomplete client
      editor/               CodeMirror setup, http language, extensions (env var hover, placeholders)
      components/           shared building blocks (QueryBar, TrafficTable, MessageViewer, SplitPane, TabStrip, SessionSidebar, …)
      ui/                   thin styled wrappers over bits-ui (Button, Input, Select, Dialog, Menu, Badge, …)
      plugins/              frontend plugin host + SDK implementation
      utils/
    pages/
      sitemap/ scopes/ filters/ intercept/ http-history/ ws-history/ match-replace/
      replay/ automate/ workflows/ assistant/ search/ findings/ exports/
      environment/ files/ plugins/ workspace/ settings/
        <Page>.svelte + local components + <page>.svelte.ts state
    styles/ tokens.css app.css
    assets/ fonts/ logo.svg
  tests/e2e/                Playwright specs
```

## 3. Svelte 5 rules
- Runes only: `$state`, `$derived`, `$effect` (sparingly), `$props`, `$bindable`. **No** legacy `export let`, `$:`, or `svelte/store` in new code.
- Shared state = classes in `*.svelte.ts` files with `$state` fields, exported as singletons or created per page context
  (`setContext`/`getContext` with typed keys). Example:
  ```ts
  // pages/replay/replay.svelte.ts
  export class ReplayState {
    sessions = $state<Session[]>([]);
    activeId = $state<number | null>(null);
    active = $derived(this.sessions.find((s) => s.id === this.activeId) ?? null);
    async send() { /* calls api, updates entries */ }
  }
  ```
- Events: callback props (`onsend`), not `createEventDispatcher`. Snippets (`{#snippet}` / `{@render}`) instead of slots.
- `$effect` only for syncing with non-Svelte things (CodeMirror, subscriptions); always return cleanup.
- Large lists never rendered without virtualization. Avoid deep `$state` proxies for big arrays of rows → use `$state.raw` and replace immutably.
- Components ≤ 250 lines; split otherwise. One component per file, PascalCase names.

## 4. Data flow
- **Server is the source of truth.** Pages query via urql; live updates via subscriptions merge into page state.
- No REST calls except files/exports/backups/CA/preview (`lib/api/rest.ts`).
- Every GraphQL operation lives in a `.graphql` file next to the page (`pages/replay/replay.graphql`), codegen produces typed documents.
- Optimistic UI only for cheap toggles (row colour, rule enabled). Errors → toast with `extensions.code` mapped to human text.
- UI preferences (column config, pane sizes, sidebar layout) persist to backend `ui_state` (per project) or `settings` (instance), not localStorage.

## 5. Commands, shortcuts, context menus
- All actions registered in `lib/commands/registry.ts`: `{ id, title, icon?, when?(ctx), run(ctx) }`.
- Shortcuts map to command ids (`lib/commands/shortcuts.ts`), loaded from backend settings, rebindable.
- Context menus are built from command ids + `when` predicates — never ad-hoc onclick menus.

## 6. Styling rules
- Only design tokens (CSS variables) for colour/spacing/radius/typography; no hard-coded hex in components.
- Dark/light via `data-theme` on `<html>`. Density via `data-density`.
- Tailwind utilities for layout; component-specific CSS in `<style>` when clearer.
- No Caido assets, colours or icon set (see `00-overview.md §4`).

## 7. Testing
- Unit/component tests next to files (`Foo.test.ts`). Test behaviour via Testing Library queries by role/label.
- E2E (Playwright): `bun run test:e2e` builds the Go binary, starts it on a random port with a temp data dir and a local target server, then drives Chromium through the proxy. Each feature milestone adds E2E specs mirroring the AC in `01-feature-spec.md`.
- Visual regressions: Playwright screenshots for shell + key pages in dark/light (threshold 0.2%).

## 8. Scripts (`web/package.json`)
```
dev        bunx --bun vite                 # proxies /api to http://127.0.0.1:8080
build      bunx --bun vite build           # outputs web/dist (embedded by Go)
check      svelte-check && eslint . && prettier --check .
test       vitest run
test:e2e   playwright test
codegen    graphql-codegen --config codegen.ts   # reads ../internal/api/graphql/schema/*.graphqls
```
