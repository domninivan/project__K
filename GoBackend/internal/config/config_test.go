package config

import (
	"strings"
	"testing"
)

func TestLoadErrors(t *testing.T) {
	base := map[string]string{
		"GO_DB_HOST": "localhost",
		"GO_DB_USER": "user",
		"GO_DB_NAME": "testdb",
		"GO_DB_PORT": "5432",
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr string
	}{
		{
			name:    "empty GO_DB_HOST",
			env:     merge(base, "GO_DB_HOST", ""),
			wantErr: "GO_DB_HOST is required",
		},
		{
			name:    "empty GO_DB_USER",
			env:     merge(base, "GO_DB_USER", ""),
			wantErr: "GO_DB_USER is required",
		},
		{
			name:    "empty GO_DB_NAME",
			env:     merge(base, "GO_DB_NAME", ""),
			wantErr: "GO_DB_NAME is required",
		},
		{
			name:    "non-numeric GO_DB_PORT",
			env:     merge(base, "GO_DB_PORT", "abc"),
			wantErr: "invalid GO_DB_PORT",
		},
		{
			name:    "zero GO_DB_PORT",
			env:     merge(base, "GO_DB_PORT", "0"),
			wantErr: "invalid GO_DB_PORT",
		},
		{
			name:    "GO_DB_PORT out of range",
			env:     merge(base, "GO_DB_PORT", "70000"),
			wantErr: "invalid GO_DB_PORT",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			cfg, err := Load()

			if err == nil {
				t.Fatalf("expected error, got nil (cfg=%+v)", cfg)
			}
			if cfg != nil {
				t.Errorf("expected nil cfg, got %+v", cfg)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantErr)
			}
		})
	}
}

func merge(base map[string]string, k, v string) map[string]string {
	out := make(map[string]string, len(base))
	for kk, vv := range base {
		out[kk] = vv
	}
	out[k] = v
	return out
}
