---
id: T-013
title: Contract - history queries, subscriptions, completions, colour/note mutations
milestone: M1
role: lead
status: todo
depends_on: [T-007]
owner: ""
branch: task/T-013-contract-history
paths:
  - internal/api/graphql/schema/history.graphqls
  - internal/api/graphql/schema/query.graphqls
  - internal/api/graphql/resolvers/stubs_history.go
  - internal/api/graphql/generated/**
  - web/src/lib/api/generated/**
  - web/src/lib/api/mocks/**
spec:
  - docs/01-feature-spec.md#F-HIST-1..4
  - docs/04-query-language.md §7
  - docs/03-architecture.md §6
---

## Goal
Schema for `requests(first, after, query, scopeId, sort)` connection, `request(id)`, `requestAdded(query, scopeId)` / `requestUpdated` subscriptions,
`queryCompletions`, `setRequestColor`, `setRequestNote`, `deleteRequests(query)`; stub resolvers; UI mocks with realistic fake data.

## Acceptance criteria
- [ ] All fields documented; Bytes fields base64; large bodies fetched via separate field with `offset/limit`.
- [ ] `make gen` clean; frontend can render mock history with no backend.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
