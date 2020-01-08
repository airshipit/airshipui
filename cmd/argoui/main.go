package main

import (
	"fmt"
	"log"

	"opendev.org/airship/airshipui/internal/environment"
	"opendev.org/airship/airshipui/internal/plugin/argoui"
)

var pluginName = "argoui"

func main() {
	// Remove the prefix from the go logger since Octant will print logs with timestamps.
	log.SetPrefix("")

	description := fmt.Sprintf("Argo UI version %s", environment.Version())
	// Use the plugin service helper to register this plugin.
	p, err := plugin.Register(pluginName, description)
	if err != nil {
		log.Fatal(err)
	}

	// The plugin can log and the log messages will show up in Octant.
	log.Printf("%s is starting", pluginName)
	p.Serve()
}
