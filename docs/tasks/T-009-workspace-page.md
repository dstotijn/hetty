---
id: T-009
title: Workspace page (projects) and top-bar project switcher
milestone: M0
role: frontend
status: todo
depends_on: [T-003, T-006, T-007]
owner: ""
branch: task/T-009-workspace-page
paths:
  - web/src/pages/workspace/**
  - web/src/lib/components/ProjectSwitcher*.svelte
  - web/src/lib/stores/project.svelte.ts
spec:
  - docs/01-feature-spec.md#F-WSP-1
  - docs/02-ui-layout.md §1, §3.12
---

## Goal
Users manage projects on the Workspace page (Projects tab; Backups tab placeholder) and switch projects from the top bar; all page stores reset on `projectChanged`.

## Acceptance criteria
- [ ] Create/rename/duplicate/delete/set default/select work against the real API (component tests with mocked client + one Playwright smoke).
- [ ] Switcher is keyboard accessible, searchable, shows the temp project distinctly.
- [ ] Deleting requires typed confirmation of the project name.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
