package extflag

import (
	"fmt"
	"strings"
)

// ToMapOfStringToStringSlice takes a string, like "1=10.0.0.1,10.0.0.2:2=10.0.0.10", and returns
// its equivalent map[string][]string object: { 1: [10.0.0.1, 10.0.0.2], 2: [10.0.0.10] }
// Map entries MUST be separated by `:`. Slice entries MUST be separated by `,`
// Slices can be null, for example "1:2:3=foo,bar" would return { "1": nil, "2": nil, "3": ['foo' 'bar'] }
func ToMapOfStringToStringSlice(val string) map[string][]string {
	return ToCustomMapOfStringToStringSlice(val, ",", ":", "=")
}

func ToCustomMapOfStringToStringSlice(val, sliceSep, mapEntrySep, assignOperator string) map[string][]string {
	if len(val) == 0 {
		return nil
	}

	parts := strings.Split(val, mapEntrySep)
	m := make(map[string][]string)
	for _, part := range parts {
		kv := strings.Split(part, assignOperator)
		switch len(kv) {
		case 2:
			m[kv[0]] = strings.Split(kv[1], sliceSep)
			break
		case 1:
			m[kv[0]] = nil
			break
		default:
			// Bogus input for this part (i.e. input is like: "1=good,input:bogus:2=im,ok,too")
			continue
		}
	}

	fmt.Printf("ToMapOfStringToStringSlice => %+v", m)

	return m
}

// KebabCaseToCamelCase converts kebab-style-strings to CAMEL_CASE_STRINGS,
// useful in binding flags to their equivalent environment variable name
func KebabCaseToCamelCase(kebab string) string {
	return strings.ToUpper(strings.ReplaceAll(kebab, "-", "_"))
}
