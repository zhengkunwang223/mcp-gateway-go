package main

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/zhengkunwang223/mcp-gateway-go/gateway"
)

func main() {
	s := server.NewMCPServer(
        "SSE Demo",
        "1.0.0",
	)
	sseServer := gateway.NewSSEServer(s,gateway.WithBaseURL("http://127.0.0.1:7979"))
	mcpCommand := gateway.McpCommand{
		Command: "npx",
		Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
	}
	gatewayServer := gateway.NewGatewaySSEServer(mcpCommand, sseServer)
	if err := gatewayServer.Start(":7979");err != nil {
		panic(err)
	}
}