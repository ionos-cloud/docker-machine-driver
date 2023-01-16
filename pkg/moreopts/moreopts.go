package moreopts

// This package focuses on processing String type returned by drivers.DriverOptions interface into more complex types

import (
	"fmt"
	"github.com/docker/machine/libmachine/drivers"
	"strings"
)

// StringToStringSlice queries the String value at key, and expects the `,` separator for new slice entries, and `:` for new map entries
// e.g. 1=10.0.0.1,10.0.0.2:2=20.1.0.10 => map[int][]string{ 1: { "10.0.0.1", "10.0.0.2" }, 2: { "20.1.0.10" } }
func StringToStringSlice(opts drivers.DriverOptions, key string) map[string][]string {
	val := opts.String(key)
	out := make(map[string][]string)
	lans := strings.Split(val, ":")
	for _, lan := range lans {
		parts := strings.Split(lan, "=")
		lanId, ips := parts[0], parts[1]
		out[lanId] = append(out[lanId], strings.Split(ips, ",")...)
	}
	fmt.Printf("Out %+v", out)
	return out
}
