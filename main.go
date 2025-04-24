package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/zhengkunwang223/mcp-gateway-go/gateway"
)

type Config struct {
	Name 		 string
	Version 	 string
	Port         int
	BaseURL      string
	SSEPath      string
	MessagePath  string
	Command 	 string	
	Headers      []string
	OAuth2Bearer string
}

func main() {
	config := Config{}
	flag.StringVar(&config.Name, "name", "mcp-gateway-go", "SSE Server name")
	flag.StringVar(&config.Version, "version", "1.0.0", "SSE Server version")
	flag.IntVar(&config.Port, "port", 7979, "Server listening port")
	flag.StringVar(&config.BaseURL, "baseUrl", "http://127.0.0.1:7979", "Base URL")
	flag.StringVar(&config.SSEPath, "ssePath", "/sse", "Base URL for SSE")
	flag.StringVar(&config.MessagePath, "messagePath", "/messages", "Path for messages")
	flag.StringVar(&config.Command, "command", "", "Command to run")
	flag.StringVar(&config.OAuth2Bearer, "oauth2Bearer", "ceshi", "OAuth2 Bearer token")

	flag.Parse()
	

	s := server.NewMCPServer(
		config.Name,
		config.Version,
	)
	var opt []gateway.SSEOption
	opt = append(opt, 
		gateway.WithBaseURL(config.BaseURL),
		gateway.WithSSEEndpoint(config.SSEPath),
		gateway.WithMessageEndpoint(config.MessagePath),
	)
	if config.OAuth2Bearer != "" {
		opt = append(opt, gateway.WithOAuth2Bearer(config.OAuth2Bearer))
	}
	sseServer := gateway.NewSSEServer(s, opt...)
	var mcpCommand gateway.McpCommand
	if config.Command != "" {
		commands := strings.Split(config.Command, " ")
		if len(commands) > 0 {
			mcpCommand = gateway.McpCommand{
				Command: commands[0],
			}
			if len(commands) > 1 {
				mcpCommand.Args = commands[1:]
			}
		}
	}
	gatewayServer := gateway.NewGatewaySSEServer(mcpCommand, sseServer)
	if err := gatewayServer.Start(fmt.Sprintf(":%d",config.Port)); err != nil {
		panic(err)
	}
}