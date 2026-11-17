# ADR-0004 — GraphQL (gqlgen) with subscriptions as the only UI API

- Status: accepted
- Date: 2026-11-17

## Decision
Schema-first gqlgen; one `.graphqls` file per domain; subscriptions over `graphql-transport-ws`. Small REST surface only
for binary transfers. Frontend uses urql + graphql-codegen (typed documents). Schema changes are "contract PRs" owned by the lead.

## Consequences
+ Matches Hetty's existing approach; typed end-to-end; live updates first-class; plugins and MCP can reuse the same services.
− Codegen on both sides → generated files must never be hand-merged; conflicts are resolved by regenerating.

## Alternatives
Connect RPC (Hetty's unmerged experiment) — good typing, but subscriptions/streaming to browsers is less ergonomic; REST+SSE — more hand-written types.
