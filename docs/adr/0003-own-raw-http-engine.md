# ADR-0003 — Own raw HTTP/1.x engine for proxy, Replay and Automate

- Status: accepted
- Date: 2026-11-17

## Context
`net/http` and `httputil.ReverseProxy` canonicalize header case/order, reject malformed messages, and rewrite
framing. A security proxy must forward and send exactly what the user wrote (request smuggling, header
injection, parser-differential testing).

## Decision
`internal/rawhttp` implements a lenient HTTP/1.x reader/writer over `net.Conn` that preserves bytes, with a
structured *view* (method, target, headers with original casing and order, body boundaries) computed lazily.
The proxy front-end, Replay, Automate and workflow "Send Request" all use it. `net/http` is used only for the UI/API server.
HTTP/2 for Replay (M7) uses `golang.org/x/net/http2` Framer + hpack directly.

## Consequences
+ Byte-exact fidelity; full control over connection reuse and last-byte synchronization.
− We own correctness of chunked encoding, keep-alive, pipelining, 1xx handling → heavy fuzz + golden tests required.
