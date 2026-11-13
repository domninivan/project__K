package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/domninivan/project__K/GoBackend/internal/handler"
	"github.com/domninivan/project__K/GoBackend/internal/model"
	"github.com/domninivan/project__K/GoBackend/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

type fakeWorkCenterService struct {
	listResp   []model.WorkCenter
	listErr    error
	createResp *model.WorkCenter
	createErr  error
}

func (f *fakeWorkCenterService) List(ctx context.Context) ([]model.WorkCenter, error) {
	return f.listResp, f.listErr
}

func (f *fakeWorkCenterService) Create(ctx context.Context, req model.CreateWorkCenterRequest) (*model.WorkCenter, error) {
	return f.createResp, f.createErr
}

type fakeHealth struct{ err error }

func (f *fakeHealth) Check(ctx context.Context) error { return f.err }

func newTestRouter(svc *fakeWorkCenterService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return router.New(router.Handlers{
		Health:     handler.NewHealth(&fakeHealth{}, log),
		WorkCenter: handler.NewWorkCenter(svc, log),
	}, log)
}

func doJSON(t *testing.T, r *gin.Engine, method, target, body string, hdr map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateWorkCenter(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcResp    *model.WorkCenter
		svcErr     error
		wantStatus int
		want       *model.WorkCenter
	}{
		{
			name:       "created",
			body:       `{"code":"WC-01","title":"Press","active":false,"operation_types":["OT-1","OT-2"]}`,
			svcResp:    &model.WorkCenter{ID: 1, Code: "WC-01", Title: "Press", Active: false, OperationTypes: []string{"OT-1", "OT-2"}},
			want:       &model.WorkCenter{ID: 1, Code: "WC-01", Title: "Press", Active: false, OperationTypes: []string{"OT-1", "OT-2"}},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "malformed json",
			body:       `{"code":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation: empty code",
			body:       `{"code":"","title":"Press","operation_types":["OT-1"]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "conflict",
			body:       `{"code":"WC-01","title":"Press","operation_types":["OT-1"]}`,
			svcErr:     fmt.Errorf("insert work_center: %w", model.ErrConflict),
			wantStatus: http.StatusConflict,
		},
		{
			name:       "internal error",
			body:       `{"code":"WC-01","title":"Press","operation_types":["OT-1"]}`,
			svcErr:     errors.New("db down"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := &fakeWorkCenterService{createResp: tc.svcResp, createErr: tc.svcErr}
			w := doJSON(t, newTestRouter(svc), http.MethodPost, "/api/v1/work-centers", tc.body, nil)

			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", w.Code, tc.wantStatus, w.Body)
			}

			if tc.want != nil {
				var got model.WorkCenter
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("unmarshal %s: %v", w.Body, err)
				}
				if diff := cmp.Diff(tc.want, &got); diff != "" {
					t.Errorf("response mismatch (-want +got):\n%s", diff)
				}
				return
			}

			var resp map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal %s: %v", w.Body, err)
			}
			if _, ok := resp["error"]; !ok {
				t.Errorf("no \"error\" key in %s", w.Body)
			}
		})
	}
}

func TestListWorkCenters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		svcResp    []model.WorkCenter
		svcErr     error
		wantStatus int
		want       []model.WorkCenter
	}{
		{
			name: "ok",
			svcResp: []model.WorkCenter{
				{ID: 1, Code: "WC-01", Title: "Press", Active: true, OperationTypes: []string{"OT-1", "OT-2"}},
			},
			want: []model.WorkCenter{
				{ID: 1, Code: "WC-01", Title: "Press", Active: true, OperationTypes: []string{"OT-1", "OT-2"}},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "internal error",
			svcErr:     errors.New("db down"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := &fakeWorkCenterService{listResp: tc.svcResp, listErr: tc.svcErr}
			w := doJSON(t, newTestRouter(svc), http.MethodGet, "/api/v1/work-centers", "", nil)
			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", w.Code, tc.wantStatus, w.Body)
			}

			if tc.want != nil {
				var got []model.WorkCenter
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("unmarshal %s: %v", w.Body, err)
				}
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("response mismatch (-want +got):\n%s", diff)
				}
				return
			}

			var resp map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal %s: %v", w.Body, err)
			}
			if _, ok := resp["error"]; !ok {
				t.Errorf("no \"error\" key in %s", w.Body)
			}

		})
	}

}

func TestRequestIDEcho(t *testing.T) {
	t.Parallel()

	svc := &fakeWorkCenterService{listResp: []model.WorkCenter{}}
	w := doJSON(t, newTestRouter(svc), http.MethodGet, "/api/v1/work-centers", "", map[string]string{
		"X-Request-Id": "test-id-123",
	})

	if got := w.Header().Get("X-Request-Id"); got != "test-id-123" {
		t.Errorf("X-Request-Id = %q, want %q", got, "test-id-123")
	}

}
