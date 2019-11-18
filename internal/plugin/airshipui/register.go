/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Register(name string, description string) (*service.Plugin, error) {
	capabilities := &plugin.Capabilities{
		IsModule: true,
	}

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithNavigation(handleNavigation, initRoutes),
	}

	// Use the plugin service helper to register this plugin.
	return service.Register(name, description, capabilities, options...)
}

// handlePrint creates a navigation tree for this plugin. Navigation is dynamic and will
// be called frequently from Octant. Navigation is a tree of `Navigation` structs.
// The plugin can use whatever paths it likes since these paths can be namespaced to the
// the plugin.
func handleNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "Airship UI",
		Path:     request.GeneratePath(),
		IconName: "cloud",
	}, nil
}

// initRoutes routes for this plugin. In this example, there is a global catch all route
// that will return the content for every single path.
func initRoutes(router *service.Router) {
	router.HandleFunc("*", func(request *service.Request) (component.ContentResponse, error) {
		contentResponse := component.NewContentResponse(component.TitleFromString("Airship UI"))

		text := component.NewText("This is the Airship UI plugin.")
		contentResponse.Add(text)

		return *contentResponse, nil
	})
}
