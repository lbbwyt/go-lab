package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewSSEMCPClient("http://localhost:12345/sse")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}
	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatal(err)
	}
	listResults, err := cli.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range listResults.Tools {
		fmt.Println(t.Name + ":" + t.Description)

	}

	// 1. 构造请求
	req := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	req.Params.Name = "calculate" // 服务器端注册的工具名
	req.Params.Arguments = map[string]any{
		"operation": "add",
		"x":         3.0,
		"y":         4.0,
	}
	resp, err := cli.CallTool(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	// 3. 取结果
	fmt.Println(resp.Content[0].(mcp.TextContent).Text)
}
