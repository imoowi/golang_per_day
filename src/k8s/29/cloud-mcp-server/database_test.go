package main

import (
	"testing"

	pgquery "github.com/pganalyze/pg_query_go/v2"
)

// TestIsReadOnlySQL 覆盖 README 第 5 节表格里声明的验收标准，
// 外加几条同属该函数错误分支的边界。
func TestIsReadOnlySQL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool // true = 允许执行
	}{
		// —— 文章表格里的 9 条 ——
		{"普通查询", "SELECT id, name FROM users ORDER BY id DESC LIMIT 10", true},
		{"CTE 里是查询", "WITH c AS (SELECT 1) SELECT * FROM c", true},
		{"多语句注入", "SELECT 1; DROP TABLE users", false},
		{"DELETE", "DELETE FROM users", false},
		{"UPDATE", "UPDATE t SET x=1", false},
		{"INSERT ... RETURNING", "INSERT INTO t (x) VALUES (1) RETURNING id", false},
		{"TRUNCATE", "TRUNCATE t", false},
		{"SELECT ... FOR UPDATE 写锁", "SELECT * FROM t FOR UPDATE", false},
		{"CTE 里藏写操作", "WITH c AS (DELETE FROM t RETURNING *) SELECT * FROM c", false},

		// —— 文章没列、但同属这个函数错误分支的边界 ——
		{"空串", "", false},
		{"纯空白", "   \n\t ", false},
		{"SELECT INTO 会建表", "SELECT * INTO new_t FROM t", false},
		{"只是以 SELECT 开头", "SELECTX FROM t", false},
		{"注释起头", "-- 注释\nSELECT 1", true},
		{"DROP TABLE", "DROP TABLE users", false},
		{"语法就不合法", "SELECT FROM WHERE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isReadOnlySQL(tt.sql)
			if tt.want && err != nil {
				t.Fatalf("isReadOnlySQL(%q) = %v; want nil", tt.sql, err)
			}
			if !tt.want && err == nil {
				t.Fatalf("isReadOnlySQL(%q) = nil; want error", tt.sql)
			}
		})
	}
}

// TestStmtName 确认拒绝时给出的语句名是人看得懂的，而不是 protobuf 的默认输出。
func TestStmtName(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{"DELETE", "DELETE FROM users", "DELETE"},
		{"UPDATE", "UPDATE t SET x=1", "UPDATE"},
		{"INSERT", "INSERT INTO t (x) VALUES (1)", "INSERT"},
		{"TRUNCATE", "TRUNCATE t", "TRUNCATE"},
		{"DROP", "DROP TABLE users", "DROP"},
		{"ALTER TABLE", "ALTER TABLE t ADD COLUMN y int", "ALTER TABLE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := pgquery.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tt.sql, err)
			}
			got := stmtName(tree.Stmts[0].Stmt)
			if got != tt.want {
				t.Fatalf("stmtName(%q) = %q; want %q", tt.sql, got, tt.want)
			}
		})
	}
}
