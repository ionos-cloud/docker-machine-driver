package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// Since the SDKs can't yet process API errors, we must do some string ops.
func extractValueForMessageKey(s string) string {
	stripped := regexp.MustCompile(`\s+`).ReplaceAllString(strings.ReplaceAll(s, "\n", " "), " ")
	if !strings.Contains(stripped, "\"message\"") {
		// Sadly, in this case, we don't know how to process the API error.
		return stripped
	}
	return strings.Split(strings.Split(stripped, "\"message\"")[1], "\"")[1]
}

func extractIncludedJson(str string) string {
	if strings.Contains(str, "{") && strings.Contains(str, "}") {
		_, json, _ := strings.Cut(str, "{")
		return "{" + json
	}
	return str
}

func SanitizeErrorJsonToHuman(jsonErr error) error {
	dst := &bytes.Buffer{}
	str := extractIncludedJson(jsonErr.Error())
	if err := json.Compact(dst, []byte(str)); err != nil {
		//Not a valid JSON. Must manually extract the message
		return errors.New(extractValueForMessageKey(str))
	}
	return errors.New(dst.String())
}
