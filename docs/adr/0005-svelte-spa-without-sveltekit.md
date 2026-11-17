# ADR-0005 — Plain Svelte 5 + Vite SPA (no SvelteKit), Bun toolchain

- Status: accepted
- Date: 2026-11-17

## Context
The UI is a desktop-like single-page app served from the Go binary; no SSR, no server routes, no form actions.
SvelteKit 3 shipped on 2026-10-01 with config/alias changes, which adds churn without benefit for this app.

## Decision
Svelte 5 (runes) mounted by Vite, a ~150-line hash router in `web/src/lib/router`, lazy-loaded pages. Bun is the only
package manager and script runner. Build output `web/dist` is embedded with `go:embed`.

## Consequences
+ Minimal framework surface; trivial embedding; no SSR pitfalls.
− We own routing (small) and code-splitting setup (Vite handles it).
