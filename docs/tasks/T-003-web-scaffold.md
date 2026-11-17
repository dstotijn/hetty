---
id: T-003
title: Svelte 5 + Vite + Bun scaffold with app shell, router and embedding
milestone: M0
role: frontend
status: todo
depends_on: [T-001]
owner: ""
branch: task/T-003-web-scaffold
paths:
  - web/**
  - internal/api/static/**
spec:
  - docs/06-frontend-conventions.md (all)
  - docs/02-ui-layout.md §1, §2, §10
  - docs/adr/0005-svelte-spa-without-sveltekit.md
---

## Goal
`web/` contains the SPA scaffold with the app shell (TopBar, Sidebar with all categories/pages as empty placeholder pages, StatusBar),
a hash router with lazy pages, design tokens (dark/light), and the Go package `internal/api/static` that embeds `web/dist` and serves it with SPA fallback.

## Scope
- In: `bun create vite` equivalent set up by hand (svelte-ts), Tailwind v4, bits-ui, @lucide/svelte, fonts (Inter, JetBrains Mono, self-hosted),
  eslint/prettier/svelte-check/vitest config, scripts from docs/06 §8, router (`#/page/…`), placeholder pages for every sidebar item,
  command palette shell (Ctrl+K opens, lists "Go to <page>" commands), theme toggle, `internal/api/static` with `go:embed` + tests.
- Out: GraphQL client (T-007), real page content, mounting the handler in the main server (T-007 does that; you expose `static.Handler() http.Handler`).

## Acceptance criteria
- [ ] `cd web && bun install && bun run build && bun run check && bun run test` all pass.
- [ ] Every sidebar item from docs/02 §1 navigates to its placeholder page; reload keeps the route.
- [ ] Ctrl/⌘+K opens the palette; typing "repl" + Enter navigates to Replay.
- [ ] Dark/light themes switch via user menu; tokens only (no hex in components).
- [ ] Go test proves `/` serves index.html and unknown paths fall back to index.html; `/api/*` is not swallowed.
- [ ] First-load JS ≤ 150 KB gzip at this stage.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
