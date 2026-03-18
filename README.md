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
cd mcp-server
go build ./cmd/server
```

This creates a `server` binary in the current directory.

## Usage

### Running the Server

The server communicates over stdio (stdin/stdout), which is the standard transport for MCP servers:

```bash
./server
```

### Configuration

The server uses a single, realistic browser user-agent for web requests:
```
Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0
```

HTTP requests have a 30-second timeout.

## License

This project is licensed under the MIT-style license.
