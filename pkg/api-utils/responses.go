package api_utils

import (
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type MapStatusCodeMessages map[int]string

// CustomStatusCodeMessages is looked up with a status code for custom responses.
// Add your custom messages here instead of creating if statements in the ionoscloud package
var CustomStatusCodeMessages = MapStatusCodeMessages{
	404: "Resource is missing",
}

func (m MapStatusCodeMessages) Has(k int) bool {
	_, ok := m[k]
	return ok
}

type fmtPrintFunc func(string, ...any)

// SanitizeResponse calls SanitizeResponseCustom with some default customized messages. Refer to its documentation for behaviour
func SanitizeResponse(response ionoscloud.APIResponse, validCodeLogFunc fmtPrintFunc) error {
	return SanitizeResponseCustom(response, CustomStatusCodeMessages, validCodeLogFunc)
}

// SanitizeResponseCustom is responsible for breaking execution if the response passed as a parameter has a bad status code (i.e. >299).
// Refer to https://developer.mozilla.org/en-US/docs/Web/HTTP/Status for types of HTTP codes.
// If a custom response is found, but the response code is valid (<300), then we log the response.
func SanitizeResponseCustom(response ionoscloud.APIResponse, mapOfCustomResponses MapStatusCodeMessages, validCodeLogFunc fmtPrintFunc) error {
	if sc := response.StatusCode; sc < 300 {
		if mapOfCustomResponses.Has(sc) {
			validCodeLogFunc(mapOfCustomResponses[sc])
		}
		return nil
	}
	return fmt.Errorf("%d: %s", response.StatusCode, response.Message)
}
