/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"fmt"

	// corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
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
		IconName: "folder",
		Children: []navigation.Navigation{
			{
				Title:    "Argo",
				Path:     request.GeneratePath("argo"),
				IconName: "cloud",
			},
		},
	}, nil
}

// initRoutes routes for this plugin. In this example, there is a global catch all route
// that will return the content for every single path.
func initRoutes(router *service.Router) {
	router.HandleFunc("", func(request *service.Request) (component.ContentResponse, error) {

		contentResponse := component.NewContentResponse(component.TitleFromString("Airship UI"))
		contentResponse.Add(component.NewText(fmt.Sprintf("This is the Airship UI plugin")))

		return *contentResponse, nil
	})

	router.HandleFunc("/argo", func(request *service.Request) (component.ContentResponse, error) {
		contentResponse := component.NewContentResponse(component.TitleFromString("Argo Workflows"))

		// Verify that argo-ui is deployed before displaying its iframe.  Octant has visibility
		// as to whether a port forward has been created, so it is possible that the iframe
		// shows an empty frame in that situation
		errMsg := "The Argo UI is not available for the currently selected context"
		key := store.Key{APIVersion: "apps/v1", Kind: "Deployment", Namespace: "argo", Name: "argo-ui"}
		_, found, err := request.DashboardClient().Get(request.Context(), key)
		if err != nil || !found {
			contentResponse.Add(component.NewText(errMsg))
		} else {
			frame := component.NewIFrame("http://127.0.0.1:8001/workflows", "Argo Workflows UI")
			contentResponse.Add(frame)
		}

		return *contentResponse, nil
	})
}
