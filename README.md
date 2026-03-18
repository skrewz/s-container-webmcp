# Web Search MCP Server

An MCP (Model Context Protocol) server that provides web search and URL fetching capabilities with markdown output.

## Features

This MCP server provides two tools:

### 1. `search` - Web Search Tool

Searches the web using DuckDuckGo and returns results formatted as markdown.

**Input:**
- `query` (string): The search query

**Output:**
- `results` (string): Markdown-formatted search results with titles, URLs, and context snippets

<details><summary>Example</summary>
<pre>
{
  "name": "search",
  "arguments": {
    "query": "golang tutorial"
  }
}
</pre>
</details>


### 2. `get_url` - URL Fetcher Tool

Fetches a URL and converts the HTML content to markdown.

**Input:**
- `url` (string): The URL to fetch

**Output:**
- `content` (string): Markdown-formatted page content

<details><summary>Example</summary>
<pre>
{
  "name": "get_url",
  "arguments": {
    "url": "https://example.com"
  }
}
</pre>
</details>

## Installation

### Prerequisites

- Go 1.25 or later

### Build from Source

```bash
make build
```

This creates a `server` binary in the current directory.

## Usage

### Running the Server


### MCP Configuration

Generally, you'll run a container next to your IDE or MCP consumer. Run `make run` to see how. You'd then refer to it from your MCP configuration. E.g. from OpenCode:

```json
{
    "$schema": "https://opencode.ai/config.json",
    // [...]
    "mcp": {
        "webmcp": {
          "type": "remote",
          "url": "http://localhost:3952/sse",
          "enabled": true,
        },
    // [...]
```

## Development

```bash
make build    # Build the server binary
make test     # Run all tests
make clean    # Remove built binary
```

## Container Support

### Building and running the Container

```bash
make container-build   # Build with Podman
make container-run     # Run on port 3952
```

### Stdio Transport in Container

To run with stdio transport (for MCP clients):

```bash
podman run -i web-search-mcp -transport stdio
```

## Licence

This project is licensed under the MIT-style license.
