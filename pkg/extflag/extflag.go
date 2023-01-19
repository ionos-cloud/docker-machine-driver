package extflag

import (
	"fmt"
	"strings"
)

func ToMapOfStringToStringSlice(val string) map[string][]string {
	out := make(map[string][]string)
	mapping := strings.Split(val, ":")
	for _, pair := range mapping {
		parts := strings.Split(pair, "=")
		key, values := parts[0], parts[1]
		out[key] = append(out[key], strings.Split(values, ",")...)
	}
	fmt.Printf("Out %+v", out)
	return out
}

// KebabCaseToCamelCase converts kebab-style-strings to CAMEL_CASE_STRINGS,
// useful in binding flags to their equivalent environment variable name
func KebabCaseToCamelCase(kebab string) string {
	return strings.ToUpper(strings.ReplaceAll(kebab, "-", "_"))
}
