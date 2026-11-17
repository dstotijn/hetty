---
id: T-007
title: GraphQL server scaffold, subscriptions, errors, auth/CSRF, codegen both sides
milestone: M0
role: lead
status: todo
depends_on: [T-004, T-003]
owner: ""
branch: task/T-007-graphql-scaffold
paths:
  - internal/api/graphql/**
  - internal/api/server.go
  - internal/api/middleware/**
  - internal/api/rest/router.go
  - internal/app/**
  - gqlgen.yml
  - web/codegen.ts
  - web/src/lib/api/**
  - web/package.json
  - web/bun.lock
  - Makefile
  - go.mod
  - go.sum
spec:
  - docs/03-architecture.md §6, §8
  - docs/adr/0004-graphql-api.md
---

## Goal
`/api/graphql` (HTTP) and `/api/graphql/ws` (graphql-transport-ws) are served; schema dir with `base.graphqls` (scalars: Time, Bytes(base64), Int64; `Node` interface; error codes),
`project.graphqls` contract; `make gen` regenerates Go and TS; the web app has an urql client with WS subscriptions and a demo `serverInfo` query shown in the status bar.

## Acceptance criteria
- [ ] `make gen` is deterministic (running twice yields no diff).
- [ ] Error mapping: service errors → `extensions.code`; query errors carry `extensions.span`.
- [ ] CSRF: requests without `X-Ratatoskr: 1` or wrong content-type are rejected; `http.CrossOriginProtection` applied.
- [ ] A `ping` subscription emits every second; the UI shows connection status in the status bar.
- [ ] Contract for projects (types, queries, mutations, `projectChanged` subscription) documented with descriptions.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
