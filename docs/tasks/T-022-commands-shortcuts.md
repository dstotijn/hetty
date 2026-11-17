---
id: T-022
title: Command registry, shortcuts manager and context menu builder
milestone: M1
role: frontend
status: todo
depends_on: [T-003]
owner: ""
branch: task/T-022-commands-shortcuts
paths:
  - web/src/lib/commands/**
  - web/src/lib/components/CommandPalette/**
  - web/src/lib/components/ContextMenu/**
spec:
  - docs/02-ui-layout.md §4, §5
  - docs/01-feature-spec.md#F-UI-1, #F-UI-2
---

## Acceptance criteria
- [ ] Commands registered with `when` predicates; palette, shortcuts and context menus read the same registry.
- [ ] Default bindings from docs/02 §4; rebinding API with conflict detection; persistence via callback (backend wiring later).
- [ ] Chorded "g then h" navigation works; shortcuts disabled inside text inputs unless marked global.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
