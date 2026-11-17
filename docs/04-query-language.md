# 04 — Query Languages (HTTPQL-compatible & StreamQL)

Package: `internal/query/httpql`, `internal/query/streamql` (shared lexer/parser core in `internal/query/core`).
Two back-ends per language: **SQL compiler** (history tables) and **in-memory evaluator** (intercept, workflows,
match & replace conditions, live subscription filtering). Both MUST give identical results — enforced by a
shared conformance test table (`testdata/conformance/*.txtar`).

## 1. Grammar (EBNF)

```ebnf
query       = [ expr ] ;
expr        = or_expr ;
or_expr     = and_expr { OR and_expr } ;
and_expr    = unary { AND unary } ;              (* AND binds tighter than OR *)
unary       = [ NOT ] primary ;                   (* NOT is our extension *)
primary     = "(" expr ")" | clause | string_lit ;
clause      = field_clause | preset_clause | source_clause ;
field_clause  = namespace "." field [ subfield ] "." operator ":" value ;
preset_clause = "preset" ":" ( string_lit | alias ) ;
source_clause = "source" ":" ( string_lit | ident ) ;
subfield    = "." ( "name" | "value" ) | "[" string_lit "]" ;   (* header.name / header.value / header["X"] *)
value       = string_lit | integer | boolean | datetime_lit ;
string_lit  = '"' { char | '\"' | '\\' } '"' ;
alias       = /[a-z0-9][a-z0-9_-]*/ ;
OR  = "OR" | "or" ;  AND = "AND" | "and" ;  NOT = "NOT" | "not" ;
comment     = "//" … EOL | "/*" … "*/" ;          (* ignored by lexer *)
```

- Bare `string_lit` at the top level or inside groups expands to `(req.raw.cont:"x" OR resp.raw.cont:"x")`.
- Integers may be written quoted or unquoted (`resp.code.eq:200` == `resp.code.eq:"200"`).
- Escapes inside strings: `\"`, `\\`, `\n`, `\r`, `\t`, `\xHH`.
- Errors carry byte offsets → UI underlines the span.

## 2. Namespaces & fields (HTTPQL-compatible)

| Field | Type | Notes |
|---|---|---|
| `req.method` | string | |
| `req.host` | string | from Host header / target |
| `req.port` | int | |
| `req.path` | string | without query |
| `req.query` | string | without leading `?` |
| `req.ext` | string | includes leading dot: `.js` |
| `req.raw` | bytes | full raw request |
| `req.body` | bytes | |
| `req.header` (+ `.name` / `.value` / `["X-Name"]`) | string | header matching (any header unless subfield) |
| `req.len` | int | bytes of line+headers+body |
| `req.tls` | bool | `true`/`false` |
| `req.created_at` | datetime | RFC3339, ISO8601, RFC2822, RFC7231, ISO9075 accepted |
| `resp.code` | int | |
| `resp.len` | int | |
| `resp.roundtrip` | int | ms |
| `resp.raw`, `resp.body` | bytes | |
| `resp.header` (+ subfields) | string | |
| `row.id` | int | |
| `preset:` | — | reference filter preset by name or alias (recursive, cycle-checked) |
| `source:` | — | `proxy` `intercept` `replay` `automate` `workflow` `plugin` (Search page) |

Our extensions (documented as such in UI autocomplete): `req.color`, `req.note`, `req.edited` (bool),
`finding.title`, `finding.reporter` (Findings page only), `NOT`.

## 3. Operators

| Op | Types | Semantics |
|---|---|---|
| `eq`, `ne` | string, int, bool | case-sensitive for strings |
| `gt`, `gte`, `lt`, `lte` | int, datetime | |
| `cont`, `ncont` | string/bytes | **case-insensitive** (ASCII fold) substring |
| `like`, `nlike` | string/bytes | SQLite LIKE semantics: `%` any run, `_` one char; ASCII case-insensitive |
| `regex`, `nregex` | string/bytes | Go RE2 syntax (no look-around/backrefs) |

Type errors (e.g. `resp.code.cont:"2"`) are parse errors, not runtime errors.

## 4. SQL compilation

- Compile to a parameterized `WHERE` fragment over `requests r LEFT JOIN responses s` plus `EXISTS` sub-selects
  for header fields (`request_headers`, `response_headers` tables). Never interpolate values into SQL text.
- `cont` → `instr(lower(col), lower(?)) > 0` for text columns; for blob columns use registered function
  `bcontains_fold(blob, ?)`. `regex` → registered deterministic function `regexp(?, col)` with a compiled-pattern LRU cache.
- `like` → `col LIKE ? ESCAPE '\'`.
- `preset:` is expanded **before** compilation (AST substitution), max depth 8.
- Benchmarks (`BenchmarkCompile*`, `BenchmarkQuery100k`) live next to the compiler; budget: simple host/code filters < 300 ms on 100k rows.

## 5. In-memory evaluator
Evaluates the AST against a `query.Message` interface (method, host, headers iterator, raw bytes, …) built from
`rawhttp` types without copying bodies. Used for intercept conditions, workflows `Matches HTTPQL`, M&R conditions,
subscription filtering.

## 6. StreamQL (WebSocket)

| Field | Type |
|---|---|
| `ws.raw` | bytes (message payload, decompressed) |
| `ws.direction` | string `client` \| `server` (extension) |
| `ws.format` | string `text` \| `binary` (extension) |
| `ws.len` | int (extension) |
| `stream.host`, `stream.path` | string |
| `stream.port` | int (extension) |
| `stream.tls` | bool (extension) |
| `row.id` | int |
| `preset:` | WS filter presets |

Same operators, logic, comments and string-literal shorthand (`"x"` → `ws.raw.cont:"x"`).

## 7. Autocomplete contract (API)
`queryCompletions(language, text, cursor)` → list of `{label, insertText, kind: namespace|field|operator|preset|value, detail}`.
Values for `req.host` / `req.method` / `resp.code` suggested from distinct values in the project (top 50).
