package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/skrewz/web-search-mcp/tools"
)

func main() {
	transport := flag.String("transport", "stdio", "Transport protocol: stdio or http")
	port := flag.Int("port", 3952, "Port to listen on for HTTP transport")
	flag.Parse()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "web-search-mcp",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Logger: slog.Default(),
		},
	)

	tools.RegisterTools(server)

	if *transport == "http" {
		addr := fmt.Sprintf(":%d", *port)
		sseHandler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
			return server
		}, nil)
		http.Handle("/sse", sseHandler)
		log.Printf("Starting web-search-mcp server on HTTP transport at %s/sse\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	} else {
		log.Println("Starting web-search-mcp server on stdio transport...")
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
