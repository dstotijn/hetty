---
id: T-006
title: Projects service and GraphQL contract
milestone: M0
role: backend
status: todo
depends_on: [T-005, T-007]
owner: ""
branch: task/T-006-projects
paths:
  - internal/project/**
  - internal/api/graphql/resolvers/project*.go
spec:
  - docs/01-feature-spec.md#F-WSP-1
  - docs/03-architecture.md §5.1
---

## Goal
Projects can be created, listed, renamed, duplicated, deleted, set as default and selected; a temporary project exists when none does;
switching emits a `projectChanged` event that other services use to reopen their DB handles.

## Notes from lead
The schema part (`project.graphqls`) is added by the lead in T-007's follow-up commit before this task starts; implement resolvers against it.

## Acceptance criteria
- [ ] All F-WSP-1 operations via GraphQL with tests (httptest + gqlgen client).
- [ ] Deleting the selected project falls back to the default or a new temp project.
- [ ] Switching closes the old DB handle (no file locks remain on Windows — CI test).
- [ ] Temp project removed on clean shutdown unless renamed/saved.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
