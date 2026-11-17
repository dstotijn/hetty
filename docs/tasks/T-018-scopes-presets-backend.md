---
id: T-018
title: Scopes and filter presets backend with preset expansion
milestone: M1
role: backend
status: todo
depends_on: [T-015]
owner: ""
branch: task/T-018-scopes-presets-backend
paths:
  - internal/scope/**
  - internal/filterpreset/**
  - internal/query/preset*.go
  - internal/api/graphql/resolvers/scope*.go
  - internal/api/graphql/resolvers/filter*.go
spec:
  - docs/01-feature-spec.md#F-SCOPE-1, #F-FILTER-1
  - docs/04-query-language.md §2, §4
---

## Notes from lead
Schema for scopes/filters will be added by the lead to `scope.graphqls`/`filter.graphqls` before work starts (same pattern as T-013).

## Acceptance criteria
- [ ] Host glob matching table tests (ports, wildcards, IDN, IPv6), deny over allow.
- [ ] `preset:` expansion with cycle detection and depth limit 8.
- [ ] Import/export JSON round-trip.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
