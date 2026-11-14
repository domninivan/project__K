package middleware

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLevelForStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		status int
		want   slog.Level
	}{
		{200, slog.LevelInfo},
		{399, slog.LevelInfo},
		{400, slog.LevelWarn},
		{404, slog.LevelWarn},
		{499, slog.LevelWarn},
		{500, slog.LevelError},
		{503, slog.LevelError},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprint(tc.status), func(t *testing.T) {
			t.Parallel()
			if got := levelForStatus(tc.status); got != tc.want {
				t.Errorf("levelForStatus(%d) = %v, want %v", tc.status, got, tc.want)
			}
		})
	}

}

func TestNewRequestID(t *testing.T) {
	t.Parallel()

	newCtx := func(header string) *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		if header != "" {
			c.Request.Header.Set(headerRequestID, header)
		}
		return c
	}

	t.Run("echoes incoming header", func(t *testing.T) {
		t.Parallel()
		if got := newRequestID(newCtx("my-id")); got != "my-id" {
			t.Errorf("got %q, want echoed %q", got, "my-id")
		}
	})

	t.Run("boundary 64 accepted, 65 regenerated", func(t *testing.T) {
		t.Parallel()
		id64 := strings.Repeat("a", 64)
		if got := newRequestID(newCtx(id64)); got != id64 {
			t.Errorf("64 chars: got %q, want echo", got)
		}
		id65 := strings.Repeat("a", 65)
		if got := newRequestID(newCtx(id65)); got == id65 {
			t.Error("65 chars: got echo, want generated")
		}
	})

	t.Run("generated is 32 hex chars", func(t *testing.T) {
		t.Parallel()
		if got := newRequestID(newCtx("")); !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(got) {
			t.Errorf("generated id %q is not 32 hex chars", got)
		}
	})

}

func TestLogging(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		wantLevel string
	}{
		{"ok", 200, "INFO"},
		{"client error", 404, "WARN"},
		{"server error", 500, "ERROR"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, nil))

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(Logging(logger))
			r.GET("/ping", func(c *gin.Context) {
				c.Status(tc.status)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping?a=1", nil))

			line := buf.String()
			for _, want := range []string{
				"level=" + tc.wantLevel,
				fmt.Sprintf("status=%d", tc.status),
				"msg=request",
				"request_id=",
			} {
				if !strings.Contains(line, want) {
					t.Errorf("log line %q: %q not found", line, want)
				}
			}
			if w.Header().Get(headerRequestID) == "" {
				t.Error("response has no X-Request-Id header")
			}

		})
	}

}
