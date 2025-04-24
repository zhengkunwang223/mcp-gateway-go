# mcp-gateway-go

`mcp-gateway-go` is a lightweight Go-based gateway that transforms a standard input/output (stdio) Model Context Protocol (MCP) server into a Server-Sent Events (SSE) endpoint. This enables real-time communication with clients such as web browsers or AI agents over HTTP.



## Installation & Usage
  
```bash

go install github.com/zhengkunwang223/mcp-gateway-go@latest

```

```bash

mcp-gateway-go --port 7979 --baseUrl http://127.0.0.1:7979  --command "npx -y @modelcontextprotocol/server-filesystem /tmp"

```

- **`--port 7979`**: Port to listen on (default: `7979`)
- **`--baseUrl "http://localhost:7979"`**: Base URL for SSE 
- **`--ssePath "/sse"`**: Path for SSE subscriptions (default: `/sse`)
- **`--messagePath "/message"`**: Path for messages (stdio→SSE default: `/message`)
- **`--oauth2Bearer "some-access-token"`**: Adds an `Authorization` header with the provided Bearer token


## 📚 Additional Resources

- [mcp-go](https://github.com/mark3labs/mcp-go)