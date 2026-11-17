---
id: T-012
title: Shared TrafficTable (virtualized) and MessageViewer components
milestone: M1
role: frontend
status: todo
depends_on: [T-003]
owner: ""
branch: task/T-012-shared-table-viewer
paths:
  - web/src/lib/components/TrafficTable/**
  - web/src/lib/components/MessageViewer/**
  - web/src/lib/components/SplitPane/**
  - web/src/lib/editor/**
  - web/package.json
  - web/bun.lock
spec:
  - docs/02-ui-layout.md §2, §3.1, §9
  - docs/01-feature-spec.md#F-HIST-1, #F-HIST-2
---

## Goal
Reusable, data-source-agnostic table and viewer components ready for History, Search, Sitemap, Automate results.

## Notes from lead
Decide `@tanstack/svelte-virtual` vs in-house virtual list by benchmark (100k rows, 60 fps scroll); record decision in Handoff.
MessageViewer modes: Raw, Pretty (JSON/XML/HTML/JS), Hex, Table (headers/params), Preview (sandboxed iframe src from REST).

## Acceptance criteria
- [ ] Storybook-like demo route `#/dev/components` (dev builds only) with 100k fake rows.
- [ ] Column show/hide/reorder/resize; sort indicators; multi-select with shift/ctrl; keyboard navigation; `onrowcontextmenu` callback.
- [ ] Viewer handles a 5 MB body without freezing (chunked read-only mode); search (Ctrl+F) works.
- [ ] Component tests for selection logic and column config persistence callbacks.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
