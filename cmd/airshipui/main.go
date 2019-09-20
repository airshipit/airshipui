package main

import (
    "log"
    "opendev.org/airship/airshipui/internal/plugin"
)

var pluginName = "airship-ui"

// This is a sample plugin showing the features of Octant's plugin API.
func main() {
    // Remove the prefix from the go logger since Octant will print logs with timestamps.
    log.SetPrefix("")

    // Use the plugin service helper to register this plugin.
    p, err := plugin.Register(pluginName, "Airship UI")
    if err != nil {
        log.Fatal(err)
    }

    // The plugin can log and the log messages will show up in Octant.
    log.Printf("%s is starting", pluginName)
    p.Serve()
}
