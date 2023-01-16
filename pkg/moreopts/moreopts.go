package moreopts

// This package focuses on processing String type returned by drivers.DriverOptions interface into more complex types

import (
	"fmt"
	"github.com/docker/machine/libmachine/drivers"
	"strconv"
	"strings"
)

// IntToStringSlice queries the String value at key, and expects the `,` separator for new slice entries, and `:` for new map entries
// e.g. 1=10.0.0.1,10.0.0.2:2=20.1.0.10 => map[int][]string{ 1: { "10.0.0.1", "10.0.0.2" }, 2: { "20.1.0.10" } }
func IntToStringSlice(opts drivers.DriverOptions, key string) (out map[int][]string) {
	val := opts.String(key)
	lans := strings.Split(val, ":")
	for _, lan := range lans {
		parts := strings.Split(lan, "=")
		strLan, ips := parts[0], parts[1]
		lanId, err := strconv.Atoi(strLan)
		if err != nil {
			panic(err)
		}
		out[lanId] = strings.Split(ips, ",")
	}
	fmt.Printf("Out %+v", out)
	return out
}
