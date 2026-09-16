# 云原生·进阶-29、如何用 Go 实现一个 MCP Server？

上一篇我们把 MCP 想明白了：它把"Agent 怎么连外部能力"这件事标准化了。

但"想明白"和"做出来"之间还隔着一步。今天不聊概念，直接动手——用 Go 写一个真正能跑的 MCP Server。

我们不写 Hello World。直接让它接上三套真实系统：

```text
AI Agent
    ↓
  MCP
    ↓
cloud-mcp-server
 ├── Kubernetes      （查 Pod）
 ├── Prometheus      （查指标）
 └── PostgreSQL      （查数据）
```

做完之后，Agent 可以直接回答这些问题：

```text
production 里现在有哪些 Pod？
user-service 最近 5 分钟的请求速率是多少？
users 表最近 10 条记录是什么？
```

先立一个规则：

> **MCP Server 不是"把 API 包一层"就完事。真正值得学的，是怎么把现有能力设计成 Agent 能理解、能发现、能安全调用的 Tool。**

今天围绕这句话，从脚手架一路写到安全边界。

---

## 一、先搭脚手架：目录与依赖

创建项目：

```bash
mkdir cloud-mcp-server && cd cloud-mcp-server
go mod init github.com/example/cloud-mcp-server
```

结构刻意把 Tool 和具体系统拆开——生产项目千万别把所有东西塞进 main.go：

```text
cloud-mcp-server/
├── main.go          # 启动、组装、注册 Tool
├── k8s.go           # Kubernetes Client + 两个 Tool
├── prometheus.go    # Prometheus HTTP 查询
├── database.go      # PostgreSQL 只读查询 + SQL 校验
└── go.mod
```

依赖三组（外加一个专用于 SQL 校验的 PG 解析器）：

```bash
go get k8s.io/client-go k8s.io/apimachinery
go get github.com/jackc/pgx/v5
go get github.com/modelcontextprotocol/go-sdk
go get github.com/pganalyze/pg_query_go/v2
```

> **macOS 上得多加一步。** `pg_query_go` 是 libpg_query 的 cgo 封装，而 cgo 构建不会跑 PostgreSQL 的 `configure`，于是它自带的那份 `strchrnul` 会和现代 macOS SDK 里的系统声明撞车，报 `static declaration of 'strchrnul' follows non-static declaration`。加个宏让它用系统自带的就好：
>
> ```bash
> export CGO_CFLAGS="-DHAVE_STRCHRNUL=1"
> ```

为什么选 Go？这套工具链本身就是 Go 的强项：Kubernetes 生态、client-go、HTTP、pgx、MCP SDK、Context、Goroutine、静态编译，全在一门语言里。

---

## 二、核心就一件事：把 Go 函数注册成 Tool

整个 demo 最关键的一步，是把一个普通 Go 函数变成 MCP Tool。先看最终代码，再拆开讲。

main.go 里，三个系统各有自己的注册函数，统一在启动时挂上：

```go
package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "cloud-mcp-server",
			Version: "0.1.0",
		},
		nil,
	)

	// 注册三组 Tool（各自返回 error，便于逐个降级）
	if err := registerK8sTools(ctx, server); err != nil {
		log.Printf("[k8s] unavailable: %v", err)
	}
	if err := registerPrometheusTool(server); err != nil {
		log.Printf("[prometheus] unavailable: %v", err)
	}
	if err := registerDatabaseTool(ctx, server); err != nil {
		log.Printf("[db] unavailable: %v", err)
	}

	// 把 server 挂到 stdio，让 Agent 进程能通过 MCP 调用
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
```

这里有两个值得留意的点。

第一个是 `mcp.NewServer(impl, nil)` 那个 `nil`。它的实际类型是 `*ServerOptions`——传 `nil` 就是用一份零值配置，SDK 内部自己补默认值。

那生产环境的鉴权、计量、依赖注入该插在哪？两个地方：一是 `ServerOptions` 上的字段（`Logger` 接 `*slog.Logger`、`GetSessionID` 给会话一个可审计的标识）；二是**中间件**——所有请求都会穿过它：

```go
// type Middleware func(MethodHandler) MethodHandler
server.AddReceivingMiddleware(authMiddleware, metricsMiddleware)
```

第二个是**分组注册、逐个降级**：哪套系统没配好就跳过哪组，别让整条链路挂掉。最后 `server.Run(ctx, &mcp.StdioTransport{})` 把 Server 挂在标准输入输出上，Agent 作为外部进程拉起来就能直接通过 MCP 协议调用。

下面逐个实现。

---

## 三、Tool 一：Kubernetes（list_pods / get_pod）

先解决"集群在哪"。生产里 MCP Server 一般跑在集群内，走 ServiceAccount 的 InClusterConfig；本地开发走 `~/.kube/config`。两个都支持：

```go
// k8s.go
func newKubernetesClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}
```

Tool 本体——注意函数签名是 `(*mcp.CallToolRequest, Input)`，这是框架注入请求和反序列化入参的约定：

```go
type ListPodsInput struct {
	Namespace string `json:"namespace"`
}

type GetPodInput struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func registerK8sTools(ctx context.Context, s *mcp.Server) error {
	client, err := newKubernetesClient()
	if err != nil {
		return err // 集群不可用时整组降级，不阻塞其他 Tool
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pods",
		Description: "List pods in a Kubernetes namespace with their phase.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListPodsInput) (*mcp.CallToolResult, any, error) {
		pods, err := client.CoreV1().Pods(in.Namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("list pods: %w", err)
		}
		out := make([]string, 0, len(pods.Items))
		for _, p := range pods.Items {
			out = append(out, fmt.Sprintf("%s [%s]", p.Name, p.Status.Phase))
		}
		return nil, out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pod",
		Description: "Get a single pod by name and namespace.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in GetPodInput) (*mcp.CallToolResult, any, error) {
		pod, err := client.CoreV1().Pods(in.Namespace).Get(ctx, in.Name, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("get pod: %w", err)
		}
		return nil, pod, nil
	})
	return nil
}
```

到这里，`Go Function → MCP Tool → Agent` 就连通了。

**为什么 Description 这么重要？** 名字 `list_pods` 本身信息有限，但 "List pods in a Kubernetes namespace with their phase" 让 LLM 清楚知道：这是干什么的、什么时候该调、要什么参数、返回什么。**Description 不是写给程序员看的，是写给 Agent 理解能力边界用的。** 这也是前面讲的 Tool Design 真正落到代码上的地方。

---

## 四、Tool 二：Prometheus（query_prometheus）

Prometheus 本身有 HTTP 查询 API，MCP Server 只是把它包一层。刻意用裸 HTTP 而不引 SDK，正好能看清"包一层"到底包了什么：

```go
// prometheus.go
// PrometheusClient 是最小可用的 Prometheus 查询客户端。
// 刻意不引官方 SDK、只用裸 HTTP——文章要看清的正是"包一层"到底包了什么。
type PrometheusClient struct {
	BaseURL string
	Client  *http.Client
}

// queryResponse 对应 Prometheus 即时查询的响应信封：
// {"status":"success","data":{"resultType":"vector","result":[...]}}
type queryResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     any    `json:"result"`
	} `json:"data"`
}

// Query 跑一次 PromQL 即时查询，返回 data.result。
func (p *PrometheusClient) Query(ctx context.Context, promQL string) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(promQL) == "" {
		return nil, nil, errors.New("empty query")
	}

	endpoint := strings.TrimRight(p.BaseURL, "/") + "/api/v1/query?" +
		url.Values{"query": {promQL}}.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("query prometheus: %w", err)
	}
	defer resp.Body.Close()

	var out queryResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	// Prometheus 出错时也会返回 200 + status=error，两个都要判
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("prometheus returned %d: %s", resp.StatusCode, out.Error)
	}
	if out.Status != "success" {
		return nil, nil, fmt.Errorf("prometheus query failed: %s", out.Error)
	}

	return nil, map[string]any{
		"resultType": out.Data.ResultType,
		"result":     out.Data.Result,
	}, nil
}

func registerPrometheusTool(s *mcp.Server) error {
	baseURL := os.Getenv("PROMETHEUS_URL")
	if baseURL == "" {
		return errors.New("PROMETHEUS_URL not set") // 没配好就整组降级，不阻塞其他 Tool
	}

	p := &PrometheusClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 15 * time.Second},
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "query_prometheus",
		Description: "Run a PromQL instant query against Prometheus. Returns vector values.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in struct {
			Query string `json:"query"`
		}) (*mcp.CallToolResult, any, error) {
			return p.Query(ctx, in.Query)
	})
	return nil
}
```

用户问"user-service 最近 5 分钟请求速率多少"，Agent 自己生成 PromQL 再调这个 Tool：

```text
rate(http_requests_total{service="user-service"}[5m])
```

**注意 Tool 名要写准。** 这个 Tool 查的是 metrics（指标），不是 logs（日志）——`query_prometheus` 这个名字比 `query_metrics` 更准确地反映了它走的是 Prometheus 的 PromQL 端点。

---

## 五、Tool 三：PostgreSQL（query_database）——安全重灾区

这是今天风险最大的一个 Tool。一个裸的 `query_database(sql)` 一旦开放，就等于把整个数据库交给了 LLM。

所以这个 Tool 有三道防线，层层收紧：

```text
① 数据库层  只读账号 + 表级权限     （最后、也是最硬的一闸）
② Tool 层   单语句 + 只读 AST 校验   （下面核心代码）
③ 运行层    行数上限 + 查询超时       （防止拖垮连接池）
```

先看 Tool 层的 SQL 校验，这是全文最该逐行读的部分。

关键选件是 `pg_query_go`——PostgreSQL 官方解析器（C 语言）的 Go 封装。**这里千万别用 `strings.HasPrefix(sql, "SELECT")` 那套字符串校验**：`"SELECT 1; DROP TABLE users"` 这种拼接注入、藏在 `WITH` 里的 `DELETE`、`SELECT ... FOR UPDATE` 行锁，字符串前缀校验统统拦不住。用解析器，这些全从 AST 上挡掉。

```go
// database.go
// go get github.com/pganalyze/pg_query_go/v2

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
```

这段校验在本地跑过 16 个用例，全部符合预期：

```text
PASS | ALLOW | SELECT id, name FROM users ORDER BY id DESC LIMIT 10
PASS | ALLOW | WITH c AS (SELECT 1) SELECT * FROM c
PASS | DENY  | SELECT 1; DROP TABLE users          （多语句注入）
PASS | DENY  | DELETE FROM users
PASS | DENY  | UPDATE t SET x=1
PASS | DENY  | INSERT ... RETURNING id
PASS | DENY  | TRUNCATE t
PASS | DENY  | SELECT * FROM t FOR UPDATE          （写锁）
PASS | DENY  | WITH c AS (DELETE ...) SELECT ...   （CTE 里的写）
```

把 SQL 净化好之后，再加行数上限和查询超时，防止一条"合法但巨耗"的 `SELECT * FROM big_table` 拖垮连接池或把 Context 撑爆：

```go
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
```

然后把这个函数注册成 Tool。注册前先用 `DATABASE_URL` 建连接池并 `Ping`，避免
`pgxpool.New` 仅仅因为 DSN 格式正确就误判数据库可用：

```go
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
```

**必须清醒的一点：** 就算用了 `pg_query_go` 解析，"能过校验"也不等于"这条查询就一定安全"——解析器能挡住多语句和写操作，但挡不住一条巨耗的 `SELECT`。所以真正的安全是**纵深防御**，上面三层之外，数据库层那个只读账号才是最硬的一闸：SQL 校验绕不过它。Demo 给到的是 Tool 层；生产环境四层（数据库 / Tool / MCP 鉴权审计 / 运行层）缺一不可。

---

## 六、把 Agent 跑起来：三种接法

Server 写好了，Agent 怎么调它？三条路。

**① 本地 MCP 客户端（Claude Code / Claude Desktop）**：配置里指向 `go build` 出的可执行文件：

```json
{
  "mcpServers": {
    "cloud-mcp": {
      "command": "/path/to/cloud-mcp-server",
      "env": {
        "PROMETHEUS_URL": "http://localhost:9090",
        "DATABASE_URL": "postgres://user:pass@localhost:5432/app"
      }
    }
  }
}
```

跑起来直接问"production 现在有哪些 Pod？"，你会看到它自动发现 `list_pods` 并调用。

**② Go MCP Client**：`for tool, err := range cs.Tools(ctx, nil) { … }` 拿到每个 Tool 的 schema，逐个调用——和 Agent 内部逻辑一致。

**③ 自写 Agent Loop**：把 `Tools()` 列出的 schema 塞进 LLM 的 tool 定义，模型输出调用意图，你执行再喂回去。就是前面那个 `Goal → Plan → Act → Observe → Re-plan`。

三种接法背后是同一套 MCP 协议——这正是它标准化的价值。

---

## 七、一次完整调用长什么样

用户问："为什么 production 的 user-service 最近响应变慢？"

Agent 自己规划、自己串 Tool，我们一行诊断逻辑都没写死：

```text
用户："为什么 production 的 user-service 最近响应变慢？"
   ↓
Planning：要看 Pod 状态 + 指标
   ↓
Tool 1: list_pods("production")
   → user-service-7f6c8d7d9c-x2abc [Running] …
   ↓
Tool 2: query_prometheus("rate(http_requests_total{service=\"user-service\"}[5m])")
   → CPU 85%，P95 从 120ms 涨到 780ms
   ↓
LLM 分析
   ↓
回答："user-service 有 4 个 Pod，过去 5 分钟 CPU 使用率 85%，
      P95 延迟从 120ms 升到 780ms，判断是 CPU 资源不足。"
```

分工非常干净：

```text
Tool 负责拿数据
LLM 负责分析数据
```

这正是 Agent Loop 和 MCP 在实战里真正合体。

---

## 八、Tool 设计原则：小而明确

最容易犯的错误是把能力做大：

```text
execute_kubectl(command)   ← 危险：等于交出整个集群
```

拆开才是对的：

```text
list_pods · get_pod · get_events · get_logs   ← 每个能力边界清晰
```

一个好 Tool 应该做到：

```text
输入明确 + 能力单一 + 输出结构化 + 权限可控 + 容易审计
```

`get_pod(namespace, name)` 的能力边界是"读一个 Pod"；`execute_kubectl` 的能力边界是"整个集群"。后者根本不该出现在生产。这和微服务"接口边界越清晰越好控制"是同一条思想。

---

## 九、MCP Server 该放在哪？

```text
个人开发：
Mac → Agent Client → MCP Server（本地进程，stdio）

生产：
Agent → MCP Gateway（鉴权 / 审计 / 限流 / 计量）
            ├── K8s MCP
            ├── Prometheus MCP
            └── DB MCP
```

到这一步，MCP 就从"开发者工具"变成了"企业 AI 基础设施"。Gateway 那层就是昨天反复强调的安全边界。

---

## 十、再理解一遍 MCP

回头看昨天的问题就清楚了：MCP 并没有让 Kubernetes、Prometheus、PostgreSQL 变强，它改变的是**这些系统如何被 Agent 发现和使用**：

```text
以前：Agent → Kubernetes SDK → Kubernetes
现在：Agent → MCP → cloud-mcp-server → Kubernetes
```

后端系统开始有了"AI 可调用接口"。过去开发者操作平台的路径是 `CLI → API`，现在逐渐变成 `Agent → MCP`——"帮我看看 production 集群有没有异常"，Agent 自己查 Pod、查事件、查指标，综合给出结论。

但别把 MCP 当魔法：

```text
MCP  ≠  Agent
MCP 负责"连接能力"，Agent 负责"目标 → 推理 → 规划 → 执行"
K8s / Prometheus / PG 各管各的领域
```

四样组合起来，才是 AI 驱动的云原生运维。

---

## 十一、给 Go 后端的启示

你以前的工作面是：写 API → 写 Service → 写 DB → 部署 K8s。

现在多了一整层：

```text
设计 AI Tool → 设计 MCP Server → 设计 Agent Runtime → 设计权限体系 → 设计 AI 可观测性
```

一句话：**你不再只给人写 API，也要开始给 Agent 设计能力接口。** 这是工程范式级的变化。

---

## 十二、今日总结

今天我们动手做了个 `cloud-mcp-server`，把 K8s / Prometheus / PostgreSQL 统一暴露成 Tool：

```text
Agent → MCP → cloud-mcp-server → { K8s, Prometheus, PostgreSQL }
```

带走四件事：

**① 一个 Go 函数 → 一个 MCP Tool，核心就 `mcp.AddTool`。**

**② Description 是写给 Agent 的，不是写给程序员看的。**

**③ 数据库 Tool 用 `pg_query_go` 解析 AST + 行数/超时，做纵深防御，而不是靠字符串前缀。**

**④ 能力要拆小，别做 `execute_kubectl`。**

核心思想始终是：

> MCP 不负责思考，也不负责规划。它负责把外部世界的能力，以标准化、可发现、可审计的方式交给 Agent。
