package main

import (
	ionoscloud "github.com/ionos-cloud/docker-machine-driver"
	"github.com/rancher/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(ionoscloud.NewDriver("", ""))
}
