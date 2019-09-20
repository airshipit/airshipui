package main

import (
    "fmt"
    "log"

    "opendev.org/airship/airshipui/internal/plugin"
)

var (
    pluginName = "airship-ui"
    // version will be overriden by ldflags supplied in Makefile
    version = "(dev-version)"
)

// This is a sample plugin showing the features of Octant's plugin API.
func main() {
    // Remove the prefix from the go logger since Octant will print logs with timestamps.
    log.SetPrefix("")

    description := fmt.Sprintf("Airship UI version %s", version)
    // Use the plugin service helper to register this plugin.
    p, err := plugin.Register(pluginName, description)
    if err != nil {
        log.Fatal(err)
    }

    // The plugin can log and the log messages will show up in Octant.
    log.Printf("%s is starting", pluginName)
    p.Serve()
}
