// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/skrewz/web-search-mcp/internal/scraper"
	"github.com/skrewz/web-search-mcp/internal/search"
)

type SearchInput struct {
	Query string `json:"query" jsonschema:"the search query"`
}

type SearchOutput struct {
	Results string `json:"results" jsonschema:"markdown-formatted search results"`
}

type GetURLInput struct {
	URL string `json:"url" jsonschema:"the URL to fetch"`
}

type GetURLOutput struct {
	Content string `json:"content" jsonschema:"markdown-formatted page content"`
}

func RegisterTools(server *mcp.Server) {
	searchTool := &mcp.Tool{
		Name:        "search",
		Description: "Search the web using DuckDuckGo and return results as markdown",
	}

	mcp.AddTool(server, searchTool, searchHandler)

	getURLTool := &mcp.Tool{
		Name:        "get_url",
		Description: "Fetch a URL and convert the HTML content to markdown (not appropriate for web searches)",
	}

	mcp.AddTool(server, getURLTool, getURLHandler)
}

func searchHandler(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchOutput, error) {
	if input.Query == "" {
		return nil, SearchOutput{}, fmt.Errorf("query cannot be empty")
	}

	searcher := search.NewSearcher(nil)
	results, err := searcher.Search(ctx, input.Query)
	if err != nil {
		return nil, SearchOutput{}, fmt.Errorf("search failed: %w", err)
	}

	markdown := search.FormatResultsMarkdown(results)

	return nil, SearchOutput{Results: markdown}, nil
}

func getURLHandler(ctx context.Context, req *mcp.CallToolRequest, input GetURLInput) (*mcp.CallToolResult, GetURLOutput, error) {
	if input.URL == "" {
		return nil, GetURLOutput{}, fmt.Errorf("URL cannot be empty")
	}

	scraper := scraper.NewScraper("")
	markdown, err := scraper.Fetch(ctx, input.URL)
	if err != nil {
		return nil, GetURLOutput{}, fmt.Errorf("fetching URL failed: %w", err)
	}

	return nil, GetURLOutput{Content: markdown}, nil
}
