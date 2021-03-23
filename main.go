package main

import (
	"github.com/rancher/machine/libmachine/drivers/plugin"

	"github.com/harvester/docker-machine-driver-harvester/harvester"
)

func main() {
	plugin.RegisterDriver(harvester.NewDriver("machine", ""))
}
