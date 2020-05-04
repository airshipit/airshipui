/*
 Copyright (c) 2020 AT&T. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package main

import (
	"log"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var pluginName = "airshipui-example-plugin"

// HelloWorldPlugin is a required struct to be an octant compliant plugin
type HelloWorldPlugin struct{}

// return a new hello world struct
func newHelloWorldPlugin() *HelloWorldPlugin {
	return &HelloWorldPlugin{}
}

// This is a sample plugin showing the features of Octant's plugin API.
func main() {
	// Remove the prefix from the go logger since Octant will print logs with timestamps.
	log.SetPrefix("")

	// Tell Octant to call this plugin when printing configuration or tabs for Pods
	capabilities := &plugin.Capabilities{
		IsModule: true,
	}

	hwp := newHelloWorldPlugin()

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithNavigation(hwp.handleNavigation, hwp.initRoutes),
	}

	// Use the plugin service helper to register this plugin.
	p, err := service.Register(pluginName, "The very smallest thing you can do", capabilities, options...)
	if err != nil {
		log.Fatal(err)
	}

	// The plugin can log and the log messages will show up in Octant.
	log.Printf("hello-world-plugin is starting")
	p.Serve()
}

// handles the navigation pane interation
func (hwp *HelloWorldPlugin) handleNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "Hello World Plugin",
		Path:     request.GeneratePath(),
		IconName: "folder",
	}, nil
}

// initRoutes routes for this plugin. In this example, there is a global catch all route
// that will return the content for every single path.
func (hwp *HelloWorldPlugin) initRoutes(router *service.Router) {
	router.HandleFunc("*", hwp.routeHandler)
}

// this function returns the octant wrapped HTML content for the page
func (hwp *HelloWorldPlugin) routeHandler(request service.Request) (component.ContentResponse, error) {
	contentResponse := component.NewContentResponse(component.TitleFromString("Hello World Title"))
	helloWorld := component.NewText("Hello World just some text on the page")
	contentResponse.Add(helloWorld)
	return *contentResponse, nil
}
