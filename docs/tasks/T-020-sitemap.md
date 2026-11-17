---
id: T-020
title: Sitemap tree (backend derivation + page)
milestone: M1
role: backend
status: todo
depends_on: [T-016]
owner: ""
branch: task/T-020-sitemap
paths:
  - internal/history/sitemap*.go
  - internal/api/graphql/resolvers/sitemap*.go
  - web/src/pages/sitemap/**
spec:
  - docs/01-feature-spec.md#F-SITE-1
  - docs/02-ui-layout.md §3.2
---

## Notes from lead
This task is split into two commits/branches if two agents are available: backend part (role backend) then UI part (role frontend). Lead adds `sitemap.graphqls` first.

## Acceptance criteria
- [ ] Tree children loaded lazily per node; counts correct; respects scope + query.
- [ ] 10k requests / 200 hosts: root load < 500 ms.
- [ ] Delete node deletes underlying requests after confirmation.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
