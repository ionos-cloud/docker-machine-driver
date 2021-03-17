package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/ionos-cloud/rancher-driver"
)

func main() {
	plugin.RegisterDriver(ionoscloud.NewDriver("", ""))
}
