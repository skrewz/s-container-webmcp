package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

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

	if *transport == "http" {
		addr := fmt.Sprintf(":%d", *port)
		streamableHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			// Extract remote port for logging
			remotePort := "unknown"
			if hostPort := r.RemoteAddr; hostPort != "" {
				if _, port, err := net.SplitHostPort(hostPort); err == nil {
					remotePort = port
				}
			}

			// Create a new server instance per request for stateless mode
			server := mcp.NewServer(
				&mcp.Implementation{
					Name:    "web-search-mcp",
					Version: "1.0.0",
				},
				&mcp.ServerOptions{
					Logger: enhancedLogger.With("remote_port", remotePort),
					InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
						baseLogger.Info("HTTP session initialized", "remote_port", remotePort)
					},
				},
			)
			tools.RegisterTools(server)
			return server
		}, &mcp.StreamableHTTPOptions{
			Stateless:    true,
			JSONResponse: true,
		})
		http.Handle("/mcp", streamableHandler)
		log.Printf("Starting web-search-mcp server on HTTP transport at %s/mcp\n", addr)
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
