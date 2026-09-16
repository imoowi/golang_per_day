package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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
