# ADR-0006 — goja as JS runtime for workflows and backend plugins

- Status: accepted (revisit at M6)
- Date: 2026-11-17

## Decision
`internal/jsrt` exposes `Runtime` / `Module` interfaces; the first implementation uses goja (pure Go) with
goja_nodejs event loop. User code may be TypeScript: transpiled with esbuild's Go API (Transform, target es2020)
before execution. Each run gets a fresh VM from a pool, CPU time limit (interrupt) and memory guard.

## Consequences
+ Pure Go, no cgo, mature, easy Go↔JS bindings.
− Not full ES2023; esbuild down-levels syntax. If plugin authors need missing built-ins, evaluate QuickJS-on-wazero behind the same interface.
