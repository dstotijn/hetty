# 05 — Go Conventions (Go 1.27)

`go.mod`: `go 1.27` (and `toolchain go1.27.x` pinned to the latest patch). CI uses the same version.
`CGO_ENABLED=0` always. Run `go fix ./...` (modernizers) before every PR — it is part of `make check`.

## 1. Modern Go we expect to see (and reviewers enforce)

Verify any API you are unsure of with `go doc <pkg>.<Symbol>` — **never guess signatures**.

### Go 1.27
| Feature | Where we use it |
|---|---|
| **Generic methods** (methods with own type params; cannot satisfy interfaces) | `store.Repo[T].ScanAll[U]`-style helpers, `events.Bus[T]` helpers, typed graph node I/O helpers. Don't force them where a function is clearer. |
| Struct literal keys may be **promoted-field selectors** | tests/fixtures with embedded structs |
| `encoding/json/v2` + `encoding/json/jsontext` (now GA; v1 runs on v2 backend) | **all new JSON code uses `json/v2`** (`json.Marshal/Unmarshal` with options, `MarshalWrite`, `UnmarshalRead`); streaming exports with `jsontext.Encoder`. Note stricter defaults (invalid UTF-8 rejected → raw bytes must be base64 or `[]byte` fields). Tag option `inline` is now `embed`. |
| New std **`uuid`** package | IDs for exports, backups, plugin tokens, MCP sessions (confirm path: `go doc uuid`). Row IDs stay SQLite integers. |
| `bytes.CutLast` / `strings.CutLast` | rawhttp parsing (e.g. last `:` in host:port) |
| `net/http` ALPN on user-provided `net.Conn` implementing `ConnectionState()` | UI server over MITM'd conns / h2 later |
| `net/http.Server.MaxHeaderValueCount` | harden UI/API server |
| `url.URL.Clone`, `url.Values.Clone` | instead of manual copies |
| `testing/synctest.Sleep` + `httptest.NewTestServer` (in-memory network) | deterministic tests of timeouts, rate limiter, intercept hold, pipeline timing logic |
| `goroutineleak` pprof profile (GA) | `--debug` endpoint; leak check in CI integration test |
| `crypto/tls`: ML-KEM hybrids, `ConnectionState.LocalCertificate`; `Config.Rand` deprecated | MITM TLS config; tests use `testing/cryptotest.SetGlobalRandom` |
| `crypto/x509` ML-DSA support, `RawSignatureAlgorithm` | CA module (we still issue ECDSA P-256 for browser compatibility) |
| `math/rand/v2` `(*Rand).N[T]` generic method | payload jitter, fuzz helpers |
| `hash/maphash.ComparableHasher` | custom LRU keys |
| `database/sql.ConvertAssign` | scanning helpers |
| `go fix` modernizers (`atomictypes`, `embedlit`, `slicesbackward`, `unsafefuncs`…), `stdversion` vet on by default | `make check` |

### Go 1.26
`new(expr)` for optional pointer fields (`Timeout: new(30*time.Second)`), self-referential generic constraints,
**`errors.AsType[T](err)`** instead of `errors.As`, `log/slog.NewMultiHandler` (stderr + ring buffer),
`testing.T.ArtifactDir` (store failing golden diffs/pcaps), `bytes.Buffer.Peek`, reflect iterator methods,
`net/http.Transport.NewClientConn`, `signal.NotifyContext` cause. `httputil.ReverseProxy.Director` is deprecated — we don't use ReverseProxy anyway.

### Go 1.22–1.25 baseline (always)
- Range-over-func iterators (`iter.Seq`, `iter.Seq2`) for streaming rows, headers, payload generators
  (`func (g Numbers) All() iter.Seq[[]byte]`), `slices`/`maps` iterator helpers, `strings.SplitSeq/Lines/FieldsSeq`.
- `for i := range n`, per-iteration loop vars (no `v := v`).
- `sync.WaitGroup.Go(fn)`; `testing/synctest`; `testing.B.Loop`; `t.Context()`.
- `os.Root` for every file access under the data dir and plugin dirs (path-traversal safe).
- `net/http` ServeMux method+wildcard patterns (`mux.HandleFunc("GET /api/rest/files/{id}", …)`).
- **`http.CrossOriginProtection`** middleware on the API.
- `unique.Make` for interning hot strings (header names, hosts) in the capture path.
- `weak` pointers for caches that must not pin memory (e.g. decoded-body cache).
- `omitzero` JSON tag where appropriate (v1 compat types).
- `tool` directives in go.mod: `go get -tool github.com/99designs/gqlgen` → run with `go tool gqlgen`.

## 2. Libraries (approved list — anything else needs an ADR)

| Purpose | Module |
|---|---|
| SQLite | `modernc.org/sqlite` |
| GraphQL server | `github.com/99designs/gqlgen` (latest), `github.com/vektah/gqlparser/v2` |
| JS runtime | `github.com/dop251/goja` (+ `goja_nodejs/eventloop`, `console`) — behind `internal/jsrt` interface |
| TS → JS transpile | `github.com/evanw/esbuild/pkg/api` (Transform only) |
| HTTP/2 framing, SOCKS5 | `golang.org/x/net/http2`, `golang.org/x/net/http2/hpack`, `golang.org/x/net/proxy` |
| Compression | `github.com/klauspost/compress` (zstd, gzip, flate), `github.com/andybalholm/brotli` |
| Crypto extras | `golang.org/x/crypto/chacha20poly1305` |
| Config / CLI | `github.com/peterbourgon/ff/v4` (flags + env + TOML) |
| Testing | stdlib, `github.com/google/go-cmp`, `golang.org/x/tools/txtar` |
| Trust store (optional CA install helper) | `github.com/smallstep/truststore` |

Always `go get <module>@latest` when adding; commit `go.sum`; run `govulncheck` (`go tool govulncheck ./...`).

## 3. Style
- `gofmt` + `go vet` + `golangci-lint` v2 (config `.golangci.yml`: errcheck, staticcheck, revive, gosec (tuned), bodyclose, contextcheck, nilerr, errorlint, perfsprint, usestdlibvars, modernize).
- Package names short, no `util`/`common`. Exported API documented. Files < 600 lines; split by concern.
- Errors: wrap with `fmt.Errorf("replay: send session %d: %w", id, err)`; sentinel errors per package (`ErrNotFound`); match with `errors.Is` / `errors.AsType`. Never `panic` outside `main`/init-impossible states.
- Context: first parameter; never stored in structs (except long-lived workers with explicit `Run(ctx)`).
- Interfaces are defined where consumed and kept small. Constructors `New(deps Deps) (*Service, error)`.
- Logging: `*slog.Logger` injected; structured keys `project_id`, `request_id`, `component`.
- No reflection-based magic in hot paths; avoid allocations in the capture loop (pool buffers with `sync.Pool`, measure).

## 4. Testing rules
- Every package has unit tests; table-driven; `t.Parallel()` where safe.
- **Golden tests** for rawhttp, query compiler (SQL text + params), M&R engine, workflow nodes — `testdata/*.txtar`, update with `-update` flag.
- **Fuzz tests**: `FuzzReadRequest`, `FuzzReadResponse`, `FuzzParseHTTPQL`, `FuzzStreamQL`, `FuzzWSFrame` (CI runs 30 s each nightly).
- **Conformance**: SQL compiler vs evaluator must agree on `testdata/conformance` corpus.
- Integration tests (`//go:build integration`) spin up the full app on random ports against `internal/testutil` servers; run with `-race`.
- Benchmarks for capture path and queries; regressions > 20% fail review.
- Use `testing/synctest` for anything time-based — **no `time.Sleep` in tests**.

## 5. Commands (Makefile targets)
```
make gen        # go generate ./... (gqlgen) + bun run codegen (web)
make check      # go fix -diff (must be empty), go vet, golangci-lint, go test -race ./..., web: bun run check && bun run test
make test-int   # integration tests
make build      # bun run build (web) → go build -trimpath -ldflags "-s -w -X main.version=…" ./cmd/ratatoskr
make fuzz       # short fuzz runs
make run        # build + ./ratatoskr --data-dir ./.dev-data
```
