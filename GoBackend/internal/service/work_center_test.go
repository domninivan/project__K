package service

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDedup(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty",
			in:   []string{},
			want: []string{},
		},
		{
			name: "single",
			in:   []string{"a"},
			want: []string{"a"},
		},
		{
			name: "mixed",
			in:   []string{"a", "b", "a", "c", "b"},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := dedup(tc.in)
			t.Parallel()
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("dedup() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
