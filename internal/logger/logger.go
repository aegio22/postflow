package logger

import (
	"context"
	"log/slog"
	"os"
)

// contextKey is the type used for context keys
type contextKey string

const loggerKey contextKey = "logger"

// NewLogger creates a new JSON structured logger
func NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// WithLogger adds the logger to the context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger from the context, or returns the default logger
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// Error logs an error with additional context
func Error(ctx context.Context, msg string, err error, args ...any) {
	logger := FromContext(ctx)
	allArgs := append([]any{"error", err}, args...)
	logger.ErrorContext(ctx, msg, allArgs...)
}

// Info logs an informational message
func Info(ctx context.Context, msg string, args ...any) {
	logger := FromContext(ctx)
	logger.InfoContext(ctx, msg, args...)
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string, args ...any) {
	logger := FromContext(ctx)
	logger.WarnContext(ctx, msg, args...)
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string, args ...any) {
	logger := FromContext(ctx)
	logger.DebugContext(ctx, msg, args...)
}
