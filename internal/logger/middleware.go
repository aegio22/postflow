package logger

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

const userIDKey contextKey = "user_id"

// Middleware returns an HTTP middleware that logs requests with structured logging
func Middleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Add logger to context
			ctx := WithLogger(r.Context(), logger)
			r = r.WithContext(ctx)

			// Wrap response writer to capture status code
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Build log attributes
			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes_written", rw.written,
				"remote_addr", r.RemoteAddr,
			}

			// Add user_id if available in context
			if userID, ok := r.Context().Value(userIDKey).(uuid.UUID); ok {
				attrs = append(attrs, "user_id", userID.String())
			}

			// Log at appropriate level based on status code
			msg := "http_request"
			if rw.statusCode >= 500 {
				logger.Error(msg, attrs...)
			} else if rw.statusCode >= 400 {
				logger.Warn(msg, attrs...)
			} else {
				logger.Info(msg, attrs...)
			}
		})
	}
}

// GetUserID extracts user ID from context for use in logging
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

// SetUserID adds user ID to context for logging
func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
