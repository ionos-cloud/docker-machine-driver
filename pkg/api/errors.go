package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func strip(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(strings.ReplaceAll(s, "\n", " "), " ")
}

func SanitizeErrorJsonToHuman(jsonErr error) error {
	dst := &bytes.Buffer{}
	err_str := jsonErr.Error()
	if err := json.Compact(dst, []byte(err_str)); err != nil {
		// Not a valid JSON, sadly
		fmt.Print("Not valid json")
		return errors.New(strip(err_str))
	}
	fmt.Print("valid json")
	return errors.New(strip(dst.String()))
}
