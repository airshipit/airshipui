/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package airshipopenstack

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
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

var ospConn *OpenstackPlugin

// NewOpenstackPlugin authenticates to Openstack and return a new openstack plugin struct
// this requires an openstack.json in the home dir of the user running the process
func NewOpenstackPlugin() *OpenstackPlugin {
	if ospConn == nil {
		opts, err := getOpenstackOpts()
		if err != nil {
			log.Printf("Error getting authentication opts %s\n", err)
		}

		client, err := openstack.AuthenticatedClient(opts)
		if err != nil {
			log.Printf("Error establishing client %s\n", err)
		}
		ospConn = &OpenstackPlugin{client}
	}
	return ospConn
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

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Printf("Error reading file %s\n", err)
	}

	var tmp optJSON
	err = json.Unmarshal(byteValue, &tmp)

	if err != nil {
		log.Printf("Error unmarshalling file %s\n", err)
	}

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
