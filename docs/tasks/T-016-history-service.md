---
id: T-016
title: History service and resolvers with live filtered subscriptions
milestone: M1
role: backend
status: todo
depends_on: [T-011, T-013, T-015]
owner: ""
branch: task/T-016-history-service
paths:
  - internal/history/**
  - internal/api/graphql/resolvers/history*.go
spec:
  - docs/01-feature-spec.md#F-HIST-1..4
  - docs/03-architecture.md §5.2, §6, §7
---

## Goal
Implement the T-013 contract: cursor pagination with sort, query+scope filtering, totals, colour/note, delete-by-query, completions (distinct values),
and `requestAdded` subscriptions filtered server-side with the evaluator.

## Acceptance criteria
- [ ] Pagination stable under concurrent inserts (cursor = sort key + id).
- [ ] Subscription filter uses evaluator; slow subscriber gets `stale` resync instead of blocking capture (test).
- [ ] Delete-by-query deletes responses/headers too (FK cascade) and returns count.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
