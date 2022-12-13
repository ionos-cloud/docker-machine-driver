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
			input:   errors.New("400 DAMN { \"Api\": \"whoopsie.. I tripped!\", \"message\": \"Someone tripped!\" }"),
			wantErr: errors.New("Someone tripped!"),
		},
		{
			name:    "Extract message - impossible (wrong error format) Here we cant do much except strip newlines and duplicated spaces",
			input:   errors.New("400 DAMN {\n \"Api\": \"whoopsie.. I tripped!\",\n \"message Someone tripped!\" }"),
			wantErr: errors.New("400 DAMN { \"Api\": \"whoopsie.. I tripped!\", \"message Someone tripped!\" }"),
		},
		{
			name:    "Valid JSON - In this case we are happiest, we can just make it compact",
			input:   errors.New("\t[\n\t{\n\t\t\"_id\": \"6398a01a9de8de1b8e577760\",\n\t\t\"message\": \"WHOOPS\"\n\t}\n\t]"),
			wantErr: errors.New("[{\"_id\":\"6398a01a9de8de1b8e577760\",\"message\":\"WHOOPS\"}]"),
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
