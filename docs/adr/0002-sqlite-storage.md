# ADR-0002 — SQLite (pure Go) per project

- Status: accepted
- Date: 2026-11-17

## Decision
`modernc.org/sqlite`, WAL mode. `instance.db` for instance-wide data, `projects/<id>/project.db` per project.
Raw messages stored as BLOBs plus parsed indexed columns and header tables; HTTPQL compiles to SQL.
Regex/byte-contains implemented as registered SQLite functions.

## Consequences
+ Single static binary, cross-compiles without cgo; projects are portable single files (easy backup/restore).
+ Query language maps naturally (LIKE semantics match HTTPQL).
− modernc is slower than cgo SQLite (~1.5–2×); mitigated by batched writes, indexes, and caching counts.

## Alternatives
bbolt/Badger (no ad-hoc queries), cgo sqlite3 (breaks static cross-builds), DuckDB (cgo, overkill), ncruces wasm sqlite (viable fallback if modernc perf is insufficient).
