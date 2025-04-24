package gateway

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)


type connection struct {
	stdin  io.Writer
	stdout io.Reader
	mutex  sync.Mutex
	reader *bufio.Reader
}

type McpCommand struct {
	Command string
	Args   []string
	Env    []string
}

type GatewaySSEServer struct {
	*SSEServer
	command McpCommand
}

func  NewGatewaySSEServer(command McpCommand, sseServer *SSEServer) *GatewaySSEServer {
	return &GatewaySSEServer{
		SSEServer: sseServer,
		command:   command,
	}
}

func (g *GatewaySSEServer) Start(addr string) error {
	cmd := exec.Command(g.command.Command, g.command.Args...)
	cmd.Env = append(os.Environ(), g.command.Env...)

	if err := g.InitStdioConn(cmd);err != nil {
		return err
	}

	return g.SSEServer.Start(addr)
}

func (s *GatewaySSEServer) InitStdioConn(cmd *exec.Cmd) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	s.SSEServer.conn = &connection{
		stdin:  stdin,
		stdout: stdout,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MCP process: %w", err)
	}

	go func() {
		cmd.Wait()
		log.Println("MCP process exited")
	}()

	return nil
}

func (s *SSEServer) forwardToStdio(message json.RawMessage) (json.RawMessage, error) {
    var msgObj map[string]interface{}
    if err := json.Unmarshal(message, &msgObj); err == nil {
        if _, hasID := msgObj["id"]; !hasID {
            s.conn.mutex.Lock()
            defer s.conn.mutex.Unlock()

            if _, err := fmt.Fprintf(s.conn.stdin, "%s\n", message); err != nil {
                return nil, fmt.Errorf("write notification failed: %w", err)
            }
            return nil, nil
        }
    }

    s.conn.mutex.Lock()
    defer s.conn.mutex.Unlock()

    if s.conn.reader == nil {
        s.conn.reader = bufio.NewReader(s.conn.stdout)
    }

    if _, err := fmt.Fprintf(s.conn.stdin, "%s\n", message); err != nil {
        return nil, fmt.Errorf("write request failed: %w", err)
    }

    line, err := s.conn.reader.ReadBytes('\n')
    if err != nil {
        return nil, fmt.Errorf("read response failed: %w", err)
    }

    return json.RawMessage(line), nil
}


func (s *SSEServer) handleMessageToStdio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSONRPCError(w, nil, mcp.INVALID_REQUEST, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		s.writeJSONRPCError(w, nil, mcp.INVALID_PARAMS, "Missing sessionId")
		return
	}
	sessionI, ok := s.sessions.Load(sessionID)
	if !ok {
		s.writeJSONRPCError(w, nil, mcp.INVALID_PARAMS, "Invalid session ID")
		return
	}
	session := sessionI.(*sseSession)

	var rawMessage json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawMessage); err != nil {
		s.writeJSONRPCError(w, nil, mcp.PARSE_ERROR, "Parse error")
		return
	}

	// Process message through MCPServer
	response, err := s.forwardToStdio(rawMessage)
	if err != nil {
		s.writeJSONRPCError(w, nil, mcp.INTERNAL_ERROR, fmt.Sprintf("MCP communication error: %v", err))
		return
	}

	// Only send response if there is one (not for notifications)
	if response != nil {
		eventData, _ := json.Marshal(response)

		// Queue the event for sending via SSE
		select {
		case session.eventQueue <- fmt.Sprintf("event: message\ndata: %s\n\n", eventData):
			// Event queued successfully
		case <-session.done:
			// Session is closed, don't try to queue
		default:
			// Queue is full, could log this
		}

		// Send HTTP response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	} else {
		// For notifications, just send 202 Accepted with no body
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.validateOAuth2Bearer(r) {
		http.Error(w, "Unauthorized: Invalid or missing Bearer token", http.StatusUnauthorized)
		return
	}
	path := r.URL.Path
	// Use exact path matching rather than Contains
	ssePath := s.CompleteSsePath()
	if ssePath != "" && path == ssePath {
		s.handleSSE(w, r)
		return
	}
	messagePath := s.CompleteMessagePath()
	if messagePath != "" && path == messagePath {
		s.handleMessageToStdio(w, r)
		return
	}

	http.NotFound(w, r)
}

func (s *SSEServer) validateOAuth2Bearer(r *http.Request) bool {
	if s.oAuth2Bearer == "" {
		return true
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return false
	}

	token := parts[1]
	return token == s.oAuth2Bearer
}

func WithOAuth2Bearer(token string) SSEOption {
	return func(s *SSEServer) {
		s.oAuth2Bearer = token
	}
}