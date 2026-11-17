package middleware

import (
	"fmt"
	"log/slog"
	mrand "math/rand/v2"
	"time"

	"github.com/gin-gonic/gin"
)

const headerRequestID = "X-Request-Id"
const ctxRequestID = "request_id"

func Logging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := newRequestID(c)
		c.Set(ctxRequestID, requestID)
		c.Writer.Header().Set(headerRequestID, requestID)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Int("size", c.Writer.Size()),
			slog.String("ip", c.ClientIP()),
			slog.Duration("duration", duration),
		}

		if q := c.Request.URL.RawQuery; q != "" {
			attrs = append(attrs, slog.String("query", q))
		}
		if ua := c.Request.UserAgent(); ua != "" {
			attrs = append(attrs, slog.String("user_agent", ua))
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		logger.LogAttrs(c.Request.Context(), levelForStatus(status), "request", attrs...)
	}
}

func newRequestID(c *gin.Context) string {
	const maxLen = 64
	id := c.GetHeader("X-Request-Id")
	if id != "" && len(id) <= maxLen {
		return id
	}
	return fmt.Sprintf("%016x%016x", mrand.Uint64(), mrand.Uint64())
}

func levelForStatus(status int) slog.Level {
	switch {
	case status >= 500:
		return slog.LevelError
	case status >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func RequestIDFrom(c *gin.Context) string {
	return c.GetString(ctxRequestID)
}
