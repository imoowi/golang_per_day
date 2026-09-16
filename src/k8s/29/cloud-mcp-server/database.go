package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	pgquery "github.com/pganalyze/pg_query_go/v2"
)

func isReadOnlySQL(sql string) error {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return errors.New("empty query")
	}

	// PG 官方解析器，一次 Parse 就给出 AST
	tree, err := pgquery.Parse(trimmed)
	if err != nil {
		// 解析失败（含多语句 "SELECT 1; DROP TABLE users"）→ 拒绝
		return fmt.Errorf("parse: %w", err)
	}

	// ParseResult.Stmts 里每条语句一个 RawStmt，条数就是分号切出来的句数。
	// 必须恰好 1 条，从根上挡住多语句注入。
	if len(tree.Stmts) != 1 {
		return fmt.Errorf("only a single statement is allowed, got %d", len(tree.Stmts))
	}

	n := tree.Stmts[0].Stmt
	if n.GetSelectStmt() == nil {
		// INSERT / UPDATE / DELETE / DROP / ALTER / TRUNCATE / SET … 一律拒绝
		return fmt.Errorf("only read-only SELECT is allowed, got %s", stmtName(n))
	}

	// 是 SELECT，还要再检查它里面有没有藏着写操作
	if s := n.GetSelectStmt(); s.IntoClause != nil || len(s.LockingClause) > 0 {
		// SELECT INTO 会建表；FOR UPDATE 会加行锁
		return errors.New("SELECT with INTO / FOR UPDATE not allowed")
	}

	// WITH 里某个 CTE 是写语句 → 整体拒绝
	if s := n.GetSelectStmt(); s.WithClause != nil {
		for _, cn := range s.WithClause.Ctes {
			if cn == nil {
				continue
			}
			ct := cn.GetCommonTableExpr()
			if ct == nil || ct.Ctequery == nil {
				continue
			}
			q := ct.Ctequery
			if q.GetInsertStmt() != nil || q.GetUpdateStmt() != nil || q.GetDeleteStmt() != nil {
				return errors.New("write CTE inside WITH not allowed")
			}
		}
	}
	return nil
}

// stmtName 把 AST 节点还原成语句名，只为了让拒绝时的错误消息说得清楚
// （"got DELETE" 比 "got %!s(*pg_query.Node=...)" 有用得多）。
// v2 没有任何现成的语句类型判定 API，只能自己对 oneof 做类型分支。
func stmtName(n *pgquery.Node) string {
	switch n.GetNode().(type) {
	case *pgquery.Node_SelectStmt:
		return "SELECT"
	case *pgquery.Node_InsertStmt:
		return "INSERT"
	case *pgquery.Node_UpdateStmt:
		return "UPDATE"
	case *pgquery.Node_DeleteStmt:
		return "DELETE"
	}

	// 其余两百来种节点（TRUNCATE / DROP / ALTER / SET …）不必逐个枚举：
	// 直接读 protobuf oneof 当前落在哪个字段上，字段名本身就是语句名，
	// 形如 delete_stmt / truncate_stmt / create_table_as_stmt。
	m := n.ProtoReflect()
	if oneof := m.Descriptor().Oneofs().ByName("node"); oneof != nil {
		if fd := m.WhichOneof(oneof); fd != nil {
			name := strings.TrimSuffix(string(fd.Name()), "_stmt")
			return strings.ToUpper(strings.ReplaceAll(name, "_", " "))
		}
	}
	return "UNKNOWN"
}

const maxQueryRows = 200 // Tool 级上限，防止 Agent 一次拉 100 万行

func queryDatabase(
	ctx context.Context,
	pool *pgxpool.Pool,
	rawSQL string,
) (any, error) {
	if err := isReadOnlySQL(rawSQL); err != nil {
		return nil, err
	}

	runCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	rows, err := pool.Query(runCtx, rawSQL)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	header := make([]string, 0, len(fields))
	for _, f := range fields {
		header = append(header, string(f.Name))
	}

	out := make([][]any, 0, 32)
	truncated := false
	for rows.Next() {
		if len(out) >= maxQueryRows {
			truncated = true // 还能取到下一行，说明是真被截断了
			break
		}
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		out = append(out, vals)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return map[string]any{
		"columns":   header,
		"rows":      out,
		"truncated": truncated,
	}, nil
}

func registerDatabaseTool(ctx context.Context, s *mcp.Server) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL not set") // 没配好就整组降级，不阻塞其他 Tool
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}

	// pgxpool.New 是惰性的：它只解析 DSN、不建连接，返回 nil 不代表库可达。
	// 所以要 Ping 一次才算真的验证过——这也是本组降级的触发点。
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close() // Ping 失败就别把池挂着
		return fmt.Errorf("ping: %w", err)
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "query_database",
		Description: "Run a read-only SELECT against PostgreSQL. Returns columns and rows.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in struct {
		SQL string `json:"sql"`
	}) (*mcp.CallToolResult, any, error) {
		result, err := queryDatabase(ctx, pool, in.SQL)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	})
	return nil
}
