package extflag

import (
	"fmt"
	"strings"
)

// ToMapOfStringToStringSlice takes a string, like "1=10.0.0.1,10.0.0.2:2=10.0.0.10", and returns
// its equivalent map[string][]string object: { 1: [10.0.0.1, 10.0.0.2], 2: [10.0.0.10] }
// Map entries MUST be separated by `:`. Slice entries MUST be separated by `,`
func ToMapOfStringToStringSlice(val string) map[string][]string {
	if len(val) == 0 || !strings.Contains(val, "=") {
		return nil
	}
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
