package api

import (
	"errors"
	"testing"
)

func TestSanitizeErrorJsonToHuman(t *testing.T) {
	tests := []struct {
		name    string
		input   error
		wantErr error
	}{
		{
			name:    "3 Newlines",
			input:   errors.New("hello\n\n\nworld"),
			wantErr: errors.New("hello world"),
		},
		{
			name:    "3 spaces",
			input:   errors.New("hello    world"),
			wantErr: errors.New("hello world"),
		},
		{
			name:    "Nothing to do",
			input:   errors.New("hello world"),
			wantErr: errors.New("hello world"),
		},
		{
			name:    "Extract message - possible",
			input:   errors.New("{}{[\"Api\":\"message\":\"Someone tripped!\"}"),
			wantErr: errors.New("Someone tripped!"),
		},
		{
			name:    "Most commonly met type of error - included JSON is selected and compacted",
			input:   errors.New("400 DAMN {\n \"Api\":\t \"whoopsie.. I tripped!\",\n\n \"message\": \"Someone tripped!\" }"),
			wantErr: errors.New("{\"Api\":\"whoopsie.. I tripped!\",\"message\":\"Someone tripped!\"}"),
		},
		{
			name:    "fallback to `message`",
			input:   errors.New("\t[\n\t\n\t\t\"_id\": \"6398a01a9de8de1b8e577760\",\n\t\t\"message\": \"WHOOPS\"\n\t}\n\t]"),
			wantErr: errors.New("WHOOPS"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeErrorJsonToHuman(tt.input); got.Error() != tt.wantErr.Error() {
				t.Errorf("SanitizeErrorJsonToHuman() have = %s, want = %s", got.Error(), tt.wantErr.Error())
			}
		})
	}
}
