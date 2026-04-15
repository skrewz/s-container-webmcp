package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/skrewz/web-search-mcp/internal/mcplogger"
	"github.com/skrewz/web-search-mcp/tools"
)

func main() {
	transport := flag.String("transport", "stdio", "Transport protocol: stdio or http")
	port := flag.Int("port", 3952, "Port to listen on for HTTP transport")
	flag.Parse()

	// Configure logging based on DEBUG env var
	var level slog.Level
	if os.Getenv("DEBUG") != "" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	baseLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	// Wrap logger to add session context and enhanced error logging
	enhancedLogger := mcplogger.NewEnhancedLogger(baseLogger.WithGroup("mcp"))

	// Counter for generating session IDs
	var sessionIDCounter atomic.Uint64

	if *transport == "http" {
		addr := fmt.Sprintf(":%d", *port)
		sseHandler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
			// Generate session ID for HTTP transport
			sid := sessionIDCounter.Add(1)
			baseLogger.Info("HTTP session starting", "session_id", sid)

			// Create a new server instance per connection to avoid race conditions
			server := mcp.NewServer(
				&mcp.Implementation{
					Name:    "web-search-mcp",
					Version: "1.0.0",
				},
				&mcp.ServerOptions{
					Logger: enhancedLogger,
					GetSessionID: func() string {
						return fmt.Sprintf("%d", sid)
					},
					InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
						sessionLogger := mcplogger.WithSession(enhancedLogger, sid)
						sessionLogger.Info("initialization completed")
					},
				},
			)
			tools.RegisterTools(server)
			return server
		}, nil)
		http.Handle("/sse", sseHandler)
		log.Printf("Starting web-search-mcp server on HTTP transport at %s/sse\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	} else {
		server := mcp.NewServer(
			&mcp.Implementation{
				Name:    "web-search-mcp",
				Version: "1.0.0",
			},
			&mcp.ServerOptions{
				Logger: enhancedLogger,
			},
		)
		tools.RegisterTools(server)
		log.Println("Starting web-search-mcp server on stdio transport...")
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
