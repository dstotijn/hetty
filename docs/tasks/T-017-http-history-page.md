---
id: T-017
title: HTTP History page with HTTPQL query bar
milestone: M1
role: frontend
status: todo
depends_on: [T-012, T-013]
owner: ""
branch: task/T-017-http-history-page
paths:
  - web/src/pages/http-history/**
  - web/src/lib/components/QueryBar/**
  - web/src/lib/query/**
spec:
  - docs/01-feature-spec.md#F-HIST-1..3
  - docs/02-ui-layout.md §3.1, §6
  - docs/04-query-language.md §7
---

## Goal
HTTP History page per layout §3.1 using TrafficTable + MessageViewer, a QueryBar with HTTPQL CodeMirror language (highlight, autocomplete via `queryCompletions`, error spans), live rows, context menu v1 (copy URL/cURL, colour, note, filter by host).

## Acceptance criteria
- [ ] Works against mocks and real backend (feature flag in dev).
- [ ] "↓ N new" pill behaviour per §6.
- [ ] Query errors shown inline without clearing the table.
- [ ] Column config persisted via `ui_state` (mock OK until backend exists).

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
