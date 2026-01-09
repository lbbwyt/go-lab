package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	startMCPServer()
	select {}
}

func startMCPServer() {
	svr := server.NewMCPServer("demo", mcp.LATEST_PROTOCOL_VERSION)
	svr.AddTool(mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arg := request.Params.Arguments.(map[string]any)
		op := arg["operation"].(string)
		x := arg["x"].(float64)
		y := arg["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultText("Cannot divide by zero"), nil
			}
			result = x / y
		}
		log.Printf("Calculated result: %.2f", result)
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				fmt.Println(e)
			}
		}()

		err := server.NewSSEServer(svr, server.WithBaseURL("http://localhost:12345")).Start("localhost:12345")

		if err != nil {
			log.Fatal(err)
		}
	}()
}
