/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// OpenstackPlugin is an example of a plugin that registers navigation and action handlers
// and performs some action against a remote API.
type OpenstackPlugin struct {
	provider *gophercloud.ProviderClient
}

// The gophercloud.AuthOptions struct ignores several key variables when attempting to unmarshal the json
// this is a workaround until a pull request can be created / accepted for gophercloud
// they may have reasons to not expose these variables, but it's not clear as to why this is the case
type optJSON struct {
	IdentityEndpoint string `json:"identityEndpoint"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	DomainID         string `json:"domainID"`
	TenantName       string `json:"tenantName"`
}

// domainCache keeps a record of the domains for use with other reporting features
// This is updated whenever we go into the getDomains compute function
var domainCache = map[string]string{}

// projectCache keeps a record of the flavors for use with other reporting features
// This is updated whenever we go into the getFlavors compute function
var flavorCache = map[string]string{}

// imageCache keeps a record of the images for use with other reporting features
// This is updated whenever we go into the getImages compute function
var imageCache = map[string]string{}

// projectCache keeps a record of the projects for use with other reporting features
// This is updated whenever we go into the getProjects identity function
var projectCache = map[string]string{}

// subnetCache keeps a record of the subnets for use with other reporting features
// This is updated whenever we go into the getSubnets network function
var subnetCache = map[string]string{}

// vmCache keeps a record of the servers for use with other reporting features
// This is updated whenever we go into the getVms compute function
var vmCache = map[string]string{}

// NewOpenstackPlugin authenticates to Openstack and return a new openstack plugin struct
// this requires an openstack.json in the home dir of the user running the process
func NewOpenstackPlugin() *OpenstackPlugin {
	opts, err := getOpenstackOpts()
	if err != nil {
		log.Printf("Error getting authentication opts %s\n", err)
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Printf("Error establishing client %s\n", err)
	}

	return &OpenstackPlugin{client}
}

// helper function to bring in the connection opts for OpenStack
func getOpenstackOpts() (gophercloud.AuthOptions, error) {
	/* The following environment variables are required for this function to work
	* - OS_USER_DOMAIN_ID
	* - OS_AUTH_URL
	* - OS_PROJECT_DOMAIN_ID
	* - OS_REGION_NAME
	* - OS_PROJECT_NAME
	* - OS_IDENTITY_API_VERSION
	* - OS_TENANT_NAME
	* - OS_TENANT_ID
	* - OS_AUTH_TYPE
	* - OS_PASSWORD
	* - OS_USERNAME
	* - OS_VOLUME_API_VERSION
	* - OS_TOKEN
	* - OS_USERID
	 */
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		// fall back on file if environment isn't setup
		opts, err = getOptsFromFile()
	}

	return opts, err
}

// fall back function if the environment isn't setup, hopefully file based is
func getOptsFromFile() (gophercloud.AuthOptions, error) {
	var fileName string
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	if runtime.GOOS == "windows" {
		fileName = home + "\\AppData\\Local\\octant\\etc\\openstack.json"
	} else {
		fileName = home + "/.config/octant/etc/openstack.json"
	}

	jsonFile, err := os.Open(fileName)
	defer jsonFile.Close()

	if err != nil {
		log.Printf("Error opening file %s\n", err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var tmp optJSON
	json.Unmarshal(byteValue, &tmp)

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: tmp.IdentityEndpoint,
		Username:         tmp.Username,
		Password:         tmp.Password,
		DomainID:         tmp.DomainID,
		TenantName:       tmp.TenantName,
		AllowReauth:      true,
	}

	return opts, err
}

// Register the plugin with octant
func Register(name string, description string) (*service.Plugin, error) {
	osp := NewOpenstackPlugin()

	// Remove the prefix from the go logger since Octant will print logs with timestamps.
	log.SetPrefix("")

	// Tell Octant to call this plugin when printing configuration or tabs for Pods
	capabilities := &plugin.Capabilities{
		IsModule: true,
	}

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithNavigation(osp.handleNavigation, osp.initRoutes),
	}

	// Use the plugin service helper to register this plugin.
	return service.Register(name, description, capabilities, options...)
}

// handlePrint creates a navigation tree for this plugin. Navigation is dynamic and will
// be called frequently from Octant. Navigation is a tree of `Navigation` structs.
// The plugin can use whatever paths it likes since these paths can be namespaced to the
// the plugin.
func (osp *OpenstackPlugin) handleNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "OpenStack",
		Path:     request.GeneratePath(),
		IconName: "cloud",
	}, nil
}

// initRoutes routes for this plugin. In this example, there is a global catch all route
// that will return the content for every single path.
func (osp *OpenstackPlugin) initRoutes(router *service.Router) {
	router.HandleFunc("", osp.routeHandler)
}

// Adds the OpenStack components to the visualization
func (osp *OpenstackPlugin) routeHandler(request service.Request) (component.ContentResponse, error) {
	response := component.NewContentResponse(component.TitleFromString("OpenStack"))
	response.Add(getDomains(osp))
	response.Add(getUsers(osp))
	response.Add(getFlavors(osp))
	response.Add(getSubnets(osp))
	response.Add(getNetworks(osp))
	response.Add(getImages(osp))
	response.Add(getProjects(osp))
	response.Add(getVMs(osp))
	response.Add(getVolumes(osp))

	return *response, nil
}
