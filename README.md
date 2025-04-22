# mcp-gateway-go

`mcp-gateway-go` is a lightweight Go-based gateway that transforms a standard input/output (stdio) Model Context Protocol (MCP) server into a Server-Sent Events (SSE) endpoint. This enables real-time communication with clients such as web browsers or AI agents over HTTP.

---

## 🚀 Overview
Built upon the [`mcp-go`](https://github.com/mark3labs/mcp-go) library, this project allows you to expose a locally running MCP server via an HTTP+SSE interface. This is particularly useful for integrating MCP servers with web-based clients or services that support SSE

---

## ⚙️ Features

- **Standard Input/Output Support**:Utilizes `mcp-go` to interact with MCP servers over stdio
- **SSE Support**:Converts MCP server output into SSE format for real-time client communication
- **Command-Line Tool Integration**:Supports running MCP servers via command-line tools like Node.js, facilitating the use of existing server implementations
- **Customizable Base URL**:Allows setting a base URL for the SSE endpoint, enabling flexible deployment configurations

---

## 🛠 Installation & Usage

### 1. Install Dependencies
Ensure Go is installed on your system. Then, retrieve the necessary package:

```bash
go get github.com/mark3labs/mcp-go
go get github.com/zhengkunwang223/mcp-gateway-go
```


### 2. Create the Gateway Server
Create a Go file (e.g., `main.go`) with the following conten:

```go
package main

import (
    "github.com/mark3labs/mcp-go/server"
    "github.com/zhengkunwang223/mcp-gateway-go/gateway"
)

func main() {
    // Initialize the MCP server
    s := server.NewMCPServer("SSE Demo", "1.0.0")

    // Create the SSE server with a custom base URL
    sseServer := gateway.NewSSEServer(s, gateway.WithBaseURL("http://127.0.0.1:7979"))

    // Define the MCP command to run
    mcpCommand := gateway.McpCommand{
        Command: "npx",
        Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
    }

    // Set up the gateway server
    gatewayServer := gateway.NewGatewaySSEServer(mcpCommand, sseServer)

    // Start the server on port 7979
    if err := gatewayServer.Start(":7979"); err != nil {
        panic(err)
    }
}
```

This code initializes an MCP server, sets up an SSE endpoint, defines the command to run the MCP server, and starts the gateway on port 797.

### 3. Run the Server
Execute the Go fil:

```bash
go run main.go
```

The server will start and listen for incoming connections on `http://127.0.0.1:7979`

---

## 📄 Example Use Cas

This setup is ideal for scenarios where you have an existing MCP server that communicates over stdio, and you want to expose it to web clients or services that support SSE. By using this gateway, you can integrate your MCP server into modern web applications without modifying the original server implementatin.

---

## 📚 Additional Resources

- [mcp-go GitHub Repository](https://github.com/mark3labs/mcp-go)

---

Feel free to explore and adapt this setup to suit your specific use cases. 