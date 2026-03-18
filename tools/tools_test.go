// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func connect(ctx context.Context, server *mcp.Server) (*mcp.ClientSession, error) {
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		return nil, err
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v1.0.0"}, nil)
	return client.Connect(ctx, t2, nil)
}

func TestSearchTool_ValidQuery(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("expected no error listing tools, got %v", err)
	}

	if len(tools.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools.Tools))
	}

	var searchTool *mcp.Tool
	for i := range tools.Tools {
		if tools.Tools[i].Name == "search" {
			searchTool = tools.Tools[i]
			break
		}
	}

	if searchTool == nil {
		t.Fatal("expected search tool to be registered")
	}

	if searchTool.Description == "" {
		t.Error("expected search tool to have a description")
	}
}

func TestSearchTool_EmptyQuery(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      "search",
		Arguments: map[string]any{"query": ""},
	}

	result, err := session.CallTool(ctx, params)
	if err == nil && !result.IsError {
		t.Fatal("expected error for empty query, got nil")
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.IsError {
		t.Error("expected result to have IsError set to true")
	}
}

func TestSearchTool_MockResults(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      "search",
		Arguments: map[string]any{"query": "test query"},
	}

	result, err := session.CallTool(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Content) == 0 {
		t.Error("expected result to have content")
	}
}

func TestGetURLTool_ValidURL(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("expected no error listing tools, got %v", err)
	}

	var getURLTool *mcp.Tool
	for i := range tools.Tools {
		if tools.Tools[i].Name == "get_url" {
			getURLTool = tools.Tools[i]
			break
		}
	}

	if getURLTool == nil {
		t.Fatal("expected get_url tool to be registered")
	}

	if getURLTool.Description == "" {
		t.Error("expected get_url tool to have a description")
	}
}

func TestGetURLTool_InvalidURL(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      "get_url",
		Arguments: map[string]any{"url": "not-a-valid-url"},
	}

	result, err := session.CallTool(ctx, params)
	if err == nil && !result.IsError {
		t.Fatal("expected error for invalid URL, got nil")
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.IsError {
		t.Error("expected result to have IsError set to true")
	}
}

func TestGetURLTool_MockHTML(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      "get_url",
		Arguments: map[string]any{"url": "data:text/html,<h1>Test</h1>"},
	}

	result, err := session.CallTool(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Content) == 0 {
		t.Error("expected result to have content")
	}
}

func TestRegisterTools(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	RegisterTools(server)

	ctx := context.Background()
	session, err := connect(ctx, server)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("expected no error listing tools, got %v", err)
	}

	if len(tools.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools.Tools))
	}

	toolNames := make(map[string]bool)
	for i := range tools.Tools {
		toolNames[tools.Tools[i].Name] = true
	}

	if !toolNames["search"] {
		t.Error("expected search tool to be registered")
	}

	if !toolNames["get_url"] {
		t.Error("expected get_url tool to be registered")
	}
}
