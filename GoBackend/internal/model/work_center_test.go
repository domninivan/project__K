package model

import (
	"testing"
)

func TestCreateWorkCenterRequest_IsActiveOrDefault(t *testing.T) {
	tr := true
	fa := false

	tests := []struct {
		name   string
		active *bool
		want   bool
	}{
		{
			name:   "nil defaults to true",
			active: nil,
			want:   true,
		},
		{
			name:   "explicit true",
			active: &tr,
			want:   true,
		},
		{
			name:   "explicit false",
			active: &fa,
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := CreateWorkCenterRequest{Active: tc.active}
			if got := r.IsActiveOrDefault(); got != tc.want {
				t.Errorf("IsActiveOrDefault() = %v, want %v", got, tc.want)
			}
		})
	}
}
