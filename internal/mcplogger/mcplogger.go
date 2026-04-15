// Package mcplogger provides enhanced logging for MCP servers with session tracking.
package mcplogger

import (
	"context"
	"log/slog"
	"runtime"
)

// EnhancedHandler wraps slog.Handler to add caller info and enhanced error context
type EnhancedHandler struct {
	handler slog.Handler
}

func (h *EnhancedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *EnhancedHandler) Handle(ctx context.Context, record slog.Record) error {
	// Add caller info for errors
	if record.Level >= slog.LevelError {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			record.AddAttrs(
				slog.String("caller_file", file),
				slog.Int("caller_line", line),
				slog.String("caller_func", funcName),
			)
		}

		// Add timestamp for error correlation
		record.AddAttrs(slog.Time("timestamp", record.Time))
	}

	return h.handler.Handle(ctx, record)
}

func (h *EnhancedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &EnhancedHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *EnhancedHandler) WithGroup(name string) slog.Handler {
	return &EnhancedHandler{handler: h.handler.WithGroup(name)}
}

// NewEnhancedLogger creates a logger with enhanced error logging
func NewEnhancedLogger(baseLogger *slog.Logger) *slog.Logger {
	return slog.New(&EnhancedHandler{handler: baseLogger.Handler()})
}

// WithSession creates a logger with session context
func WithSession(baseLogger *slog.Logger, sessionID uint64) *slog.Logger {
	return baseLogger.WithGroup("session").With("session_id", sessionID)
}
