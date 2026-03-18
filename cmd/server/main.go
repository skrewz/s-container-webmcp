// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/skrewz/web-search-mcp/tools"
)

func main() {
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

	log.Println("Starting web-search-mcp server on stdio transport...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
