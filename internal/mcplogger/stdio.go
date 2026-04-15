// Package mcplogger provides enhanced logging for MCP servers with session tracking.
// NOTE: The SDK's StdioTransport is used directly. Session tracking is done
// via the InitializedHandler callback which provides the session ID.
package mcplogger

// This file is kept for API compatibility but not used.
// The SDK handles stdio transport with its own ioConn that doesn't support session IDs.
// Session tracking is achieved through the InitializedHandler callback.
