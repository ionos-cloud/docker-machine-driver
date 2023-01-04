package sdk_utils

import (
	"errors"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortenErrSDK(t *testing.T) {
	tests := []struct {
		name    string
		args    error
		wantErr error
	}{
		{
			name:    "Not an OpenAPI generic error",
			args:    errors.New("404 Not Found {\n \"messages\" : [ {\n    \"errorCode\" : \"100\",\n    \"message\" : \"[(root).properties.ram] RAM of requested server is not a multiple of 256\"\n  } ]\n}\n}"),
			wantErr: errors.New("404 Not Found { \"messages\" : [ { \"errorCode\" : \"100\", \"message\" : \"[(root).properties.ram] RAM of requested server is not a multiple of 256\" } ]}}"),
		},
		{
			name:    "Usual SDK Error",
			args:    ionoscloud.NewGenericOpenAPIError("Damn, what a shame!", nil, nil, 100),
			wantErr: errors.New("Damn, what a shame!"),
		},
		{
			name:    "Valid JSON, Not OpenAPI generic error, Dirty",
			args:    errors.New("{ \n\n\"key\": \n[\"value\", 0.5, \n\t{ \"test\": 56, \n\t\"test2\": [true, null] }\n\t]\n}"),
			wantErr: errors.New("{ \"key\": [\"value\", 0.5, { \"test\": 56, \"test2\": [true, null] } ]}"),
		},
		{
			name:    "dont panic",
			args:    errors.New("720\n how did I get here?"),
			wantErr: errors.New("720 how did I get here?"),
		},
	}

	//assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ShortenOpenApiErr(tt.args)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeResponseCustom(t *testing.T) {
	type args struct {
		statusCode int
		msg        string
		msgs       MapStatusCodeMessages
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "sample custom err",
			args: args{
				statusCode: 404,
				msg:        "Not Found",
				msgs:       CustomStatusCodeMessages,
			},
			wantErr: errors.New("404: resource is missing: Not Found"),
		},
		{
			name: "sample 401",
			args: args{
				statusCode: 401,
				msg:        "Auth Fail",
				msgs:       CustomStatusCodeMessages,
			},
			wantErr: errors.New("401: authentication failed: Auth Fail"),
		},
		{
			name: "no custom msg",
			args: args{
				statusCode: 555,
				msg:        "server can't handle your request",
				msgs:       CustomStatusCodeMessages,
			},
			wantErr: errors.New("555: server can't handle your request"),
		},
		{
			name: "sample ok",
			args: args{
				statusCode: 200,
				msg:        "all good",
				msgs:       CustomStatusCodeMessages,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeStatusCodeCustom(tt.args.statusCode, tt.args.msg, tt.args.msgs)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
