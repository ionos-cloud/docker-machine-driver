package api

import (
	"errors"
	"testing"
)

func TestSanitizeErrorJsonToHuman(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "3 Newlines",
			err:     errors.New("hello\n\n\nworld"),
			wantErr: errors.New("hello world"),
		},
		{
			name:    "No newline",
			err:     errors.New("hello world"),
			wantErr: errors.New("hello world"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SanitizeErrorJsonToHuman(tt.err); errors.Is(err, tt.wantErr) {
				t.Errorf("SanitizeErrorJsonToHuman() have = %v, want = %v", err, tt.wantErr)
			}
		})
	}
}
