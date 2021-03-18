package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/ionos-cloud/docker-machine-driver"
)

func main() {
	plugin.RegisterDriver(ionoscloud.NewDriver("", ""))
}
