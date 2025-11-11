package main

import (
	"context"
	"log"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type PingInput struct {
	Message string `json:"message" jsonschema:"message to echo"`
}

type PingOutput struct {
	Response string `json:"response" jsonschema:"reply to echo"`
}

func Ping(ctx context.Context, req *mcp.CallToolRequest, input PingInput) (*mcp.CallToolResult, PingOutput, error) {
	if input.Message == "" {
		input.Message = "pong"
	}
	return nil, PingOutput{Response: "u said" + input.Message}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "memory-server", Version: "0.0.1"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Description: "echo a message",
	}, Ping)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal("Server failed to run", err)
	}
}
