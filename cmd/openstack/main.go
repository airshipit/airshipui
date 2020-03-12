package main

import (
	"fmt"
	"log"

	"opendev.org/airship/airshipui/internal/environment"
	plugin "opendev.org/airship/airshipui/internal/plugin/openstack"
)

var pluginName = "openstack"

func main() {
	// Remove the prefix from the go logger since Octant will print logs with timestamps.
	log.SetPrefix("")

	description := fmt.Sprintf("%s version %s", pluginName, environment.Version())

	// Use the plugin service helper to register this plugin.
	plugin.Register(pluginName, description)

	// The plugin can log and the log messages will show up in Octant.
	log.Printf("%s is starting", pluginName)
}
