# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`cloud-mcp-server` is a Go **MCP server over stdio** that exposes cloud-ops capabilities — Kubernetes, Prometheus, and read-only PostgreSQL — as tools an agent can discover and call. `README.md` is the published 每日一Go article for this project (*云原生·进阶-29、如何用 Go 实现一个 MCP Server？*), and the `.go` files mirror its code blocks.

The git root is `~/dev/imoowi/golang_per_day`, not this directory; see that repo's `src/k8s/README.md` for the article ↔ folder index.

**The source and the article are one artifact.** The four `.go` files are the article's excerpts, so an edit to them silently makes the published article wrong, and vice versa. Treat them as a pair: if you change a code block's subject, update the prose section that explains it (and the reverse).

## Building: cgo is mandatory, and macOS needs a flag

`go build ./...` compiles today, but **`CGO_ENABLED=0` cannot build this project** — `pg_query_go` wraps libpg_query's C sources. On macOS you must also set:

```bash
export CGO_CFLAGS="-DHAVE_STRCHRNUL=1"
```

Without it the cgo build of `pg_query_go` fails with `static declaration of 'strchrnul' follows non-static declaration`. The cause: `parser/src_port_snprintf.c` guards its own `strchrnul` with `#ifndef HAVE_STRCHRNUL`, but cgo never runs PostgreSQL's `configure`, so the macro stays unset while the modern macOS SDK declares `strchrnul` in `<string.h>`. That file's `#else` branch already expects the platform to supply it, so defining the macro is the intended path rather than a hack. README §1 documents this for readers.

Consequences worth remembering: the first build compiles ~470 C files (slow once, then cached); there is no `CGO_ENABLED=0` or scratch-image build; and upstream CI only ever tested Linux, so `darwin/arm64` is unverified territory.

Also note `go build ./...` writes a ~48 MB `cloud-mcp-server` binary into this directory, because `./...` resolves to a single main package. Use `go build -o /dev/null ./...` for a pure check.

The glue the article omitted is now implemented: `PrometheusClient.Query`, `registerDatabaseTool`, `stmtName`, and the `metav1` import. `go mod tidy` corrected `pg_query_go`'s `// indirect` marking and pinned `github.com/jackc/pgx/v5 v5.11.0`. The module path stays the article's `github.com/example/cloud-mcp-server` placeholder. `src/k8s/29/` is still untracked by git.

## Commands

```bash
export CGO_CFLAGS="-DHAVE_STRCHRNUL=1"   # macOS only; required for every build
go build -o /dev/null ./...              # type-check, without dropping a 48 MB binary
go vet ./...
gofmt -l .                               # currently clean
go test ./...                            # 22 subtests, all passing
go test -run TestIsReadOnlySQL ./...     # a single test
go run .                                 # starts the server — see below
```

`database_test.go` covers `isReadOnlySQL` and `stmtName`. The `isReadOnlySQL` table is the acceptance criteria from README §5 — all nine cases the article tabulates behave exactly as documented, verified rather than assumed. The article says 13 cases were run but tabulates only 9; the committed table has 16, adding edge cases (`SELECT INTO`, `SELECTX`, leading comment, `DROP`, whitespace-only) the article didn't list.

**`go run .` is not an interactive CLI.** It speaks JSON-RPC on stdin/stdout, so a bare run blocks waiting for a client — that's not a hang. To exercise it, register it with an MCP client (README §6):

```json
{"mcpServers": {"cloud-mcp": {"command": "/path/to/cloud-mcp-server",
  "env": {"PROMETHEUS_URL": "http://localhost:9090",
          "DATABASE_URL": "postgres://user:pass@localhost:5432/app"}}}}
```

Or drive it from Go: `for tool, err := range cs.Tools(ctx, nil)` yields each tool's schema for direct invocation.

## Architecture

`main.go` builds one `mcp.Server` and registers **three independent tool groups**, each returning `error`. A group that fails to initialize only logs — it is skipped, not fatal — so a missing cluster, unreachable Prometheus, or bad DSN degrades exactly that group's tools. Order: k8s → prometheus → database. `mcp.NewServer(impl, nil)` passes `nil` for the second parameter, which is `*ServerOptions` — a zero-value config, **not** the "custom registrar" the article originally claimed (that was corrected). Production auth/metering belongs on `ServerOptions` fields or in middleware added with `AddReceivingMiddleware`.

**`k8s.go`** — `newKubernetesClient` tries `rest.InClusterConfig()` (ServiceAccount, for in-cluster deployment), then falls back to `~/.kube/config` via `clientcmd.BuildConfigFromFlags`, so one binary serves both. Tools `list_pods` and `get_pod`. Tool descriptions are written for an LLM to decide *when* to call, not for a human reader — the article makes that point explicitly, so don't degrade them to terse labels.

**`prometheus.go`** — `query_prometheus` issues a PromQL instant query over plain HTTP (no SDK, deliberately, so the "wrapping" is visible), endpoint from `PROMETHEUS_URL`, 15s timeout. The agent generates the PromQL itself.

**`database.go`** — the substantive file, and the article's security centerpiece. It implements the middle of a three-layer defense the README names: a DB-layer read-only account (the hard gate, not in this repo), this Tool layer, and a runtime layer.

- `isReadOnlySQL` parses with **`pganalyze/pg_query_go`, the real PostgreSQL parser, and inspects the AST** — never string prefixes. It rejects unparseable input, anything but exactly one statement (`len(tree.Stmts) != 1`, which is what stops `SELECT 1; DROP TABLE users`), any non-`SelectStmt`, `SELECT ... INTO`, `FOR UPDATE` locking clauses, and writes hidden in `WITH` CTEs. The article is emphatic that `strings.HasPrefix(sql, "SELECT")` cannot stop any of these; keep the parser-based approach and don't let it erode into pattern matching.
- `stmtName` renders a rejected node as `"DELETE"` / `"TRUNCATE"` for the error message. v2 exposes no statement-type API, so it type-switches the four common write types and falls back to reading the protobuf `node` oneof's field name (`delete_stmt` → `DELETE`), covering all ~200 node types without enumerating them. Don't reach for `Node.String()` — that dumps protobuf text.
- `queryDatabase` adds the runtime layer: 15s context timeout, `maxQueryRows` (200) cap, and a `{columns, rows, truncated}` result so the agent can tell a complete answer from a clipped one. `truncated` is set only after probing one more `rows.Next()`, so an exactly-200-row result isn't misreported.
- `registerDatabaseTool` builds the pool from `DATABASE_URL`, then **`Ping`s it**: `pgxpool.New` is lazy and succeeds on a merely well-formed DSN, so the Ping is the real reachability probe and the group's degradation trigger (2s timeout, then `Close()` on failure).
- `database_test.go` holds the guard's acceptance table.

Adding a tool group: a new `<name>.go` exporting `register<Name>Tool(s *mcp.Server) error`, a call in `main.go` that logs on failure, and input structs carrying `json` tags — `mcp.AddTool` derives the tool's JSON schema from them.

## Conventions

Code comments are in Chinese, matching the article and its audience; keep new ones in Chinese. The README ends by previewing Day 30 (RBAC, tool permissions, human-in-the-loop), so this project is the foundation that the next article narrows — expect its tool surface to be the thing Day 30 constrains.
