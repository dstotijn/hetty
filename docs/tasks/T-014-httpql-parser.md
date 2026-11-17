---
id: T-014
title: HTTPQL-compatible lexer, parser, AST and conformance corpus
milestone: M1
role: backend
status: todo
depends_on: [T-001]
owner: ""
branch: task/T-014-httpql-parser
paths:
  - internal/query/**
spec:
  - docs/04-query-language.md §1–§3, §6
---

## Goal
`internal/query/core` + `httpql` parse queries into a typed AST with precise error spans; `streamql` reuses the core with its own field table.

## Acceptance criteria
- [ ] Grammar of docs/04 §1 fully implemented incl. comments, escapes, string shorthand, `preset:`/`source:`, `header["X"]`, `NOT`.
- [ ] Type checking of operator vs field type at parse time with helpful messages ("`cont` is not valid for integer field resp.code; use eq/gt/lt…").
- [ ] Precedence: `A AND B OR C` == `(A AND B) OR C` (test).
- [ ] `FuzzParseHTTPQL` 60 s clean; `String()` of AST re-parses to an equal AST (property test).
- [ ] Conformance corpus `testdata/conformance/*.txtar` with ≥ 80 cases (query, sample message, expected bool) — consumed by T-015.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
