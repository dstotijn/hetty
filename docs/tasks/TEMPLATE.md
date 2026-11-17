---
id: T-###
title: <imperative title>
milestone: M#
role: lead | backend | frontend | qa
status: todo            # todo | in_progress | in_review | changes_requested | blocked | done
depends_on: []          # e.g. [T-005, T-010]
owner: ""               # e.g. claude-code/opus, opencode/<model>
branch: task/T-###-slug
paths:                  # files/dirs this task may modify (globs)
  - internal/example/**
spec:                   # sections to read (keep short)
  - docs/01-feature-spec.md#F-XXX-1
  - docs/03-architecture.md §5
---

## Goal
<1–3 sentences: what exists after this task that didn't before>

## Scope
- In: …
- Out (explicitly not in this task): …

## Acceptance criteria
- [ ] AC1 — observable, testable statement
- [ ] AC2 — …
- [ ] `make check` green

## Notes from lead
<constraints, pitfalls, interfaces to use, perf budget>

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)
- What changed:
- How to verify:
- Known limitations:
- Follow-ups (proposed tasks):

## Review (reviewer)
- Verdict: approve | changes_requested
- Findings:
  - [ ] [blocker|major|minor|nit] path:line — description
