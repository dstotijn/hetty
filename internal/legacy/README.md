# internal/legacy — reference only, do not import

Frozen copies of Hetty packages, kept so the rewrite can read how the old proxy did
something (queue semantics, search syntax, MITM plumbing) without keeping the code on
any live path. They are excluded from `golangci-lint`, they are not wired into
`cmd/ratatoskr`, and they are deleted as soon as the matching `internal/` package lands.

**Nothing outside this tree may import it.** `importguard_test.go` fails the build if a
non-legacy package does (test imports included).

| Package | Fate (docs/03-architecture.md §3) |
|---|---|
| `proxy`, `proxy/intercept` | rewrite → `internal/proxy` + `internal/rawhttp`, `internal/intercept` |
| `filter` | replace → `internal/query` |
| `scope` | port → `internal/scope` |
| `chrome` | port → `internal/browser` (M7) |
| `reqlog`, `sender` | rewrite → `internal/history`, `internal/replay` |
| `proj` | rewrite → `internal/project` |
| `api` | rewrite → `internal/api/graphql` (schema owned by the lead) |
| `log` | replace → `log/slog` |

Deleted outright in T-002: `admin/` (Next.js UI, replaced by `web/`), `cmd/hetty`,
`pkg/db/bolt` (bbolt → SQLite, ADR-0002) and the two tests that only ran against the
bolt repository (`reqlog/reqlog_test.go`, `sender/sender_test.go`).
