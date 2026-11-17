---
id: T-010
title: rawhttp - lenient byte-preserving HTTP/1.x reader and writer
milestone: M1
role: backend
status: todo
depends_on: [T-001]
owner: ""
branch: task/T-010-rawhttp
paths:
  - internal/rawhttp/**
spec:
  - docs/01-feature-spec.md#F-NET-3
  - docs/adr/0003-own-raw-http-engine.md
  - docs/03-architecture.md §4
---

## Goal
`internal/rawhttp` reads requests/responses from a `*bufio.Reader` preserving exact bytes, exposes a lazy structured view
(method, target, version, headers in original order/case incl. duplicates, body framing), and writes messages byte-exact.

## Notes from lead
- Framing rules: RFC 9112 for well-formed input, but never "fix" input; if framing is ambiguous (CL+TE, duplicate CL),
  follow the same choice as `net/http` *for reading on the proxy side* and record `framing_warning`.
- Body handling: streaming reader with capture cap; chunked decoder exposes both raw and decoded views.
- Responses to HEAD and 1xx/204/304 have no body. Support `Connection: close` and keep-alive reuse; pipelining read-ahead tolerated.
- API sketch: `ReadRequest(br, Limits) (*Request, error)`, `(*Request).Raw() []byte`, `Headers() iter.Seq2[int, Header]`,
  `Body() io.Reader`, `WriteTo(w)`; `Builder` for editing that re-serializes only changed parts.

## Acceptance criteria
- [ ] ≥ 30 golden txtar cases round-trip byte-exact (malformed + normal), incl. obs-fold, bare LF, lowercase method, absolute-form, `*` target, huge headers (limit error).
- [ ] `FuzzReadRequest` / `FuzzReadResponse` run 60 s without panics; corpus committed.
- [ ] Zero allocations per header iteration in benchmark (after initial parse).
- [ ] No use of `net/http` types in the public API.

## Plan (implementer, ≤15 lines)

## Out-of-scope change needed (if any)

## Handoff (implementer)

## Review (reviewer)
