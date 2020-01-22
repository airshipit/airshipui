/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"encoding/json"
	"fmt"
	"net/url"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	v1 "k8s.io/api/core/v1"
)

// Register registers the plugin with Octant
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
		Title:    "Argo UI",
		Path:     request.GeneratePath(),
		IconName: "cloud",
	}, nil
}

// initRoutes routes for this plugin. In this example, there is a global catch all route
// that will return the content for every single path.
func initRoutes(router *service.Router) {
	router.HandleFunc("", func(request *service.Request) (component.ContentResponse, error) {
		response := component.NewContentResponse(component.TitleFromString("Argo UI"))

		u, err := getArgoUIURL(request)
		if err != nil || u == nil {
			errMsg := "The Argo UI is not available."
			response.Add(component.NewText(errMsg))
		} else {
			response.Add(component.NewIFrame(u.String(), "Argo UI"))
		}

		return *response, nil
	})
}

func getArgoUIURL(request *service.Request) (u *url.URL, err error) {
	ctx := request.Context()
	client := request.DashboardClient()

	found := false

	// client.Get is avoided here because when the plugin is first launched the
	// key is usually not present, and octant will display an error message in the
	// grpc library about marshaling nil.  client.List does not raise any such errors
	// when the key is not yet present
	ul, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Endpoints",
		Namespace:  "argo",
	})

	var data unstructured.Unstructured
	if err == nil {
		for _, item := range ul.Items {
			if item.GetName() == "argo-ui" {
				found = true
				data = item
				break
			}
		}
	}
	if !found {
		return u, err
	}

	var endpoints v1.Endpoints
	m, err := data.MarshalJSON()
	err = json.Unmarshal(m, &endpoints)
	if err != nil {
		return u, err
	}

	var addr string
	var port int32

	/*	Move through the subsets, address, and port arrays to construct a url
		More information on the structure can be found in the Kubernetes document below
		https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#endpointsubset-v1-core
	*/
	for _, s := range endpoints.Subsets {
		for _, a := range s.Addresses {
			addr = a.IP
			break
		}

		for _, p := range s.Ports {
			port = p.Port
			break
		}

		if addr != "" || port != 0 {
			break
		}
	}
	if addr == "" || port == 0 {
		return u, err
	}

	u = &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", addr, port),
	}

	return u, nil
}
