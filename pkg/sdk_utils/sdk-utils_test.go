package sdk_utils

import (
	"errors"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"
	"net/http"
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
			err := ShortenErrSDK(tt.args)
			if tt.wantErr != nil {
				assert.Equal(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeResponseCustom(t *testing.T) {
	type args struct {
		response             *ionoscloud.APIResponse
		mapOfCustomResponses MapStatusCodeMessages
		validCodeLogFunc     func(...any)
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "sample custom err",
			args: args{
				response:             &ionoscloud.APIResponse{Message: "Something went truly wrong", Response: &http.Response{StatusCode: 404}},
				mapOfCustomResponses: CustomStatusCodeMessages,
				validCodeLogFunc:     nil,
			},
			wantErr: errors.New("404: resource is missing: Something went truly wrong"),
		},
		{
			name: "sample OK",
			args: args{
				response:             &ionoscloud.APIResponse{Message: "OK", Response: &http.Response{StatusCode: 202}},
				mapOfCustomResponses: CustomStatusCodeMessages,
				validCodeLogFunc:     nil,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeResponseCustom(tt.args.response, tt.args.mapOfCustomResponses, tt.args.validCodeLogFunc)
			if tt.wantErr != nil {
				assert.Equal(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
