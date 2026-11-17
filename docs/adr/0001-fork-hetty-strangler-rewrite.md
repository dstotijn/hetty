# ADR-0001 — Fork Hetty, then rewrite incrementally ("strangler")

- Status: accepted
- Date: 2026-11-17

## Context
Hetty (MIT) provides a working MITM proxy, CA handling, intercept queue, scope and a GraphQL API (~6.2k hand-written
Go lines + ~9k generated). It has been dormant since v0.7.0 (2022). Its frontend (Next.js/React/MUI) does not fit our
Svelte requirement, its storage (bbolt) does not fit query-heavy features, and its proxy uses `httputil.ReverseProxy`,
which normalizes requests.

## Decision
Fork Hetty to keep history and attribution, but move all code into a new `internal/` layout (see `03-architecture.md`).
Hetty packages are ported or rewritten one by one; the app must build and run at the end of every milestone.
The old `admin/` and `pkg/db/bolt` are deleted in M0.

## Consequences
+ License/attribution clean; working reference code from day one.
− Most code will be replaced; we must not let old abstractions leak into the new design.

## Alternatives
Greenfield repo (cleaner, but loses the reference implementation and attribution trail).
