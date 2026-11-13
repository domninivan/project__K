package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecovery(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery(logger))
	r.GET("/boom", func(c *gin.Context) { panic("boom") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	line := buf.String()
	for _, want := range []string{"panic recovered", "boom", "stack"} {
		if !strings.Contains(line, want) {
			t.Errorf("log %q: %q not found", line, want)
		}
	}

}
