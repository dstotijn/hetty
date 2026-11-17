---
id: T-015
title: HTTPQL SQL compiler, in-memory evaluator and benchmarks
milestone: M1
role: backend
status: todo
depends_on: [T-014, T-005]
owner: ""
branch: task/T-015-httpql-sql-eval
paths:
  - internal/query/**
  - internal/store/functions*.go
spec:
  - docs/04-query-language.md §4, §5
---

## Goal
Compile ASTs to parameterized SQL over the project schema and evaluate them in memory; both pass the shared conformance corpus.

## Acceptance criteria
- [ ] Conformance corpus passes for both back-ends (same results).
- [ ] No value ever appears in SQL text (test asserts placeholders only).
- [ ] `BenchmarkQuery100k`: `req.host.eq` and `resp.code.gte` < 300 ms; `req.raw.cont` < 2 s on CI hardware (document numbers).
- [ ] Regex cache bounded (LRU 256).

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
