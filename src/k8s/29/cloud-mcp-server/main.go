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
